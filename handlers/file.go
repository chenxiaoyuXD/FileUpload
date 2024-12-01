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
	// Upload file to minio.
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Encrypt file before uploading.
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	encryptedData, err := utils.Encrypt(buf.Bytes())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	bucketName := "uploads"
	objectName := header.Filename

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, _ := minioClient.BucketExists(ctx, bucketName); !exists {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	_, err = minioClient.PutObject(ctx, bucketName, objectName, bytes.NewReader(encryptedData), int64(len(encryptedData)), minio.PutObjectOptions{})
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
		//c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/octet-stream", decryptedData)
}
