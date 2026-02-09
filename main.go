package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Port          = getEnv("PORT", "8080")
	UploadPath    = getEnv("UPLOAD_PATH", "./uploads")
	MaxUploadSize = getEnvInt("MAX_UPLOAD_SIZE_MB", 500) << 20
	StartTime     time.Time
)

type FileInfo struct {
	Name    string    `json:"name"`
	Size    string    `json:"size"`
	Bytes   int64     `json:"bytes"`
	URL     string    `json:"url"`
	ModTime time.Time `json:"mod_time"`
}

type ServerStats struct {
	Uptime      string `json:"uptime"`
	GoVersion   string `json:"go_version"`
	FileCount   int    `json:"file_count"`
	TotalSize   string `json:"total_size"`
	Platform    string `json:"platform"`
	NumCPU      int    `json:"num_cpu"`
}

func main() {
	StartTime = time.Now()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("🚀 AirGo Engine Ignited...")

	if err := os.MkdirAll(UploadPath, 0755); err != nil {
		log.Fatalf("Critical Error: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), CORSMiddleware(), MaxBodySizeMiddleware(int64(MaxUploadSize)))

	r.Static("/uploads", UploadPath)

	// --- ROUTES ---

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	r.GET("/", func(c *gin.Context) {
		c.File("./templates/index.html")
	})

	// System Stats for the Dashboard
	r.GET("/api/stats", func(c *gin.Context) {
		entries, _ := os.ReadDir(UploadPath)
		var totalSize int64
		count := 0
		for _, e := range entries {
			if !e.IsDir() {
				info, _ := e.Info()
				totalSize += info.Size()
				count++
			}
		}

		c.JSON(http.StatusOK, ServerStats{
			Uptime:    time.Since(StartTime).Round(time.Second).String(),
			GoVersion: runtime.Version(),
			FileCount: count,
			TotalSize: formatSize(totalSize),
			Platform:  runtime.GOOS,
			NumCPU:    runtime.NumCPU(),
		})
	})

	r.GET("/api/files", func(c *gin.Context) {
		entries, err := os.ReadDir(UploadPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage error"})
			return
		}

		files := make([]FileInfo, 0)
		for _, entry := range entries {
			if !entry.IsDir() {
				info, _ := entry.Info()
				files = append(files, FileInfo{
					Name:    entry.Name(),
					Size:    formatSize(info.Size()),
					Bytes:   info.Size(),
					URL:     fmt.Sprintf("/uploads/%s", entry.Name()),
					ModTime: info.ModTime(),
				})
			}
		}
		c.JSON(http.StatusOK, files)
	})

	r.POST("/api/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload"})
			return
		}

		filename := filepath.Base(file.Filename)
		dst := filepath.Join(UploadPath, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Save failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success", "file": filename})
	})

	// Delete Endpoint (New Feature)
	r.DELETE("/api/files/:name", func(c *gin.Context) {
		name := filepath.Base(c.Param("name"))
		target := filepath.Join(UploadPath, name)

		if err := os.Remove(target); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
	})

	// --- SERVER ---
	srv := &http.Server{
		Addr: "0.0.0.0:" + Port,
		Handler: r,
		ReadTimeout: 1 * time.Hour, // Long timeout for huge files
		WriteTimeout: 1 * time.Hour,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func MaxBodySizeMiddleware(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists { return value }
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil { return i }
	}
	return fallback
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit { return fmt.Sprintf("%d B", bytes) }
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
