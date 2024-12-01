package main

import (
	"Taurus_File_Upload/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.InitMinio()

	r := gin.Default()

	r.POST("/upload", handlers.UploadFile)
	r.GET("/download/:filename", handlers.DownloadFile)

	r.Run(":8080")
}
