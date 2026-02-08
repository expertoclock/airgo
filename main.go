package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Configuration with Defaults
var (
	Port          = getEnv("PORT", "8080")
	UploadPath    = getEnv("UPLOAD_PATH", "./uploads")
	MaxUploadSize = getEnvInt("MAX_UPLOAD_SIZE_MB", 500) << 20 // MB to Bytes
)

// FileInfo represents the metadata for a file
type FileInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
	URL  string `json:"url"`
}

func main() {
	// 1. Logger Setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting AirGo Service...")

	// 2. Directory Validation
	if err := os.MkdirAll(UploadPath, 0755); err != nil {
		log.Fatalf("Critical: Failed to create uploads directory: %v", err)
	}

	// 3. Gin Setup
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Middleware: Security Limits
	// Enforce body size limit BEFORE parsing multipart forms to prevent DoS
	r.Use(MaxBodySizeMiddleware(int64(MaxUploadSize)))

	// Middleware: Robust CORS
	r.Use(CORSMiddleware())

	// 4. Routes
	// Serve static files (uploads)
	r.Static("/uploads", UploadPath)

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up", "timestamp": time.Now().Unix()})
	})

	// Serve Frontend
	r.GET("/", func(c *gin.Context) {
		// Verify template exists to prevent panic
		if _, err := os.Stat("./templates/index.html"); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Frontend missing. Please redeploy.")
			return
		}
		c.File("./templates/index.html")
	})

	// API: Upload
	r.POST("/api/upload", func(c *gin.Context) {
		// 1. Limit Multipart Memory (RAM usage)
		// 8MB in RAM, rest on Disk. This prevents RAM exhaustion.
		if err := c.Request.ParseMultipartForm(8 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large or malformed request"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file part in request"})
			return
		}

		// 2. Path Traversal Protection
		filename := filepath.Base(file.Filename)
		
		// 3. Prevent overwriting hidden files or system files if any
		if filename == "." || filename == ".." || filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
			return
		}

		dst := filepath.Join(UploadPath, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Printf("Error saving file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file to disk"})
			return
		}

		log.Printf("Success: %s (%s)", filename, formatSize(file.Size))
		c.JSON(http.StatusOK, gin.H{
			"message": "Upload successful",
			"file":    filename,
			"size":    formatSize(file.Size),
		})
	})

	// API: List Files
	r.GET("/api/files", func(c *gin.Context) {
		entries, err := os.ReadDir(UploadPath)
		if err != nil {
			// If folder is missing (unlikely due to init), try to create it or return error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read storage"})
			return
		}

		files := make([]FileInfo, 0)
		for _, entry := range entries {
			if !entry.IsDir() {
				info, err := entry.Info()
				if err != nil {
					continue
				}
				files = append(files, FileInfo{
					Name: entry.Name(),
					Size: formatSize(info.Size()),
					URL:  fmt.Sprintf("/uploads/%s", entry.Name()),
				})
			}
		}
		c.JSON(http.StatusOK, files)
	})

	// 5. Server Config
	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: r,
		// Timeouts prevent Slowloris attacks
		ReadTimeout:  30 * time.Second, // Increased for large file uploads over slow LAN
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		// Limit header size
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 6. Start & Graceful Shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen failed: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 5 second timeout for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited successfully")
}

// ---------------- Helpers & Middleware ----------------

// MaxBodySizeMiddleware limits the size of the request body (payload)
func MaxBodySizeMiddleware(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing dynamically
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// If no origin, it's a direct request (like curl), allow it or ignore
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin) // Reflect Origin
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}