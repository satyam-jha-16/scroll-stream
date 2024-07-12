package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	bucketName      string
	projectID       string
	tempUploadDir   string
	tempConvertDir  string
	gcsBucketFolder string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	bucketName = os.Getenv("GCS_BUCKET_NAME")
	projectID = os.Getenv("GCS_PROJECT_ID")
	tempUploadDir = os.Getenv("TEMP_UPLOAD_DIR")
	tempConvertDir = os.Getenv("TEMP_CONVERT_DIR")
	gcsBucketFolder = os.Getenv("GCS_BUCKET_FOLDER")

	// Set GCS credentials
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("GCS_CREDENTIALS_FILE"))
}

func main() {
	router := setupRouter()
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173", "*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	router.POST("/upload", handleUpload)
	return router
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoID := uuid.New().String()
	tempFilePath := filepath.Join(tempUploadDir, videoID+filepath.Ext(file.Filename))

	if err := os.MkdirAll(tempUploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	outputPath := filepath.Join(tempConvertDir, videoID)
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output directory"})
		return
	}

	hlsPath := filepath.Join(outputPath, "index.m3u8")

	if err := convertToHLS(tempFilePath, outputPath, hlsPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert video"})
		return
	}

	gcsURL, err := uploadToGCS(outputPath, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to GCS"})
		return
	}

	// Clean up local files
	os.RemoveAll(tempFilePath)
	os.RemoveAll(outputPath)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Video converted to HLS and uploaded to GCS",
		"videoUrl": gcsURL,
		"videoId":  videoID,
	})
}

func convertToHLS(inputPath, outputPath, hlsPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-codec:v", "libx264", "-codec:a", "aac",
		"-hls_time", "10", "-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputPath, "segment%03d.ts"),
		"-start_number", "0", hlsPath)

	return cmd.Run()
}

func uploadToGCS(localDir, videoID string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		gcsPath := filepath.Join(gcsBucketFolder, videoID, relPath)
		obj := bucket.Object(gcsPath)
		writer := obj.NewWriter(ctx)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = io.Copy(writer, file); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload files: %v", err)
	}

	gcsURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s/index.m3u8", bucketName, gcsBucketFolder, videoID)
	return gcsURL, nil
}
