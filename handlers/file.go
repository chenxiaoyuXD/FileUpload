package handlers

import (
	"Taurus_File_Upload/utils"
	"bytes"
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

const defaultChunkSize = 1 * 1024 * 1024 // 1MB

func InitMinio() {
	// Initialize minio client object.
	var err error
	minioClient, err = minio.New("127.0.0.1:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
}

func UploadFile(c *gin.Context) {
	bucketName := "uploads"

	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, _ := minioClient.BucketExists(ctx, bucketName); !exists {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	// Get the file size
	fileSize := file.Size

	uploadID, err := minioClient.NewMultipartUpload(ctx, bucketName, file.Filename, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var partNumber int
	var uploadedParts []minio.CompletePart
	for offset := int64(0); offset < fileSize; offset += int64(defaultChunkSize) {
		partNumber++
		chunkSize := int64(defaultChunkSize)
		if offset+chunkSize > fileSize {
			chunkSize = fileSize - offset
		}

		// Read a chunk of data
		chunk := make([]byte, chunkSize)
		_, err := io.ReadFull(src, chunk)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Encrypt the chunk
		encryptedChunk, err := utils.Encrypt(chunk)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Upload the encrypted chunk
		uploadPart, err := minioClient.PutObject(ctx, bucketName, file.Filename, uploadID, partNumber, bytes.NewReader(encryptedChunk), int64(len(encryptedChunk)), "", "", nil)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		uploadedParts = append(uploadedParts, minio.CompletePart{
			PartNumber: partNumber,
			ETag:       uploadPart.ETag,
		})
	}

	// Complete multipart upload
	_, err = minioClient.CompleteMultipartUpload(context.Background(), bucketName, file.Filename, uploadID, uploadedParts)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}

func DownloadFile(c *gin.Context) {
	objectName := c.Param("filename")
	bucketName := "uploads"

	obj, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer obj.Close()

	// Decrypt file before downloading.
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	decryptedData, err := utils.Decrypt(buf.Bytes())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/octet-stream", decryptedData)
}
