package api

import (
	"fmt"
	"net/http"

	"airdb.io/airdb/noah/internal/version"
	"github.com/gin-gonic/gin"
)

func DefaultRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"deploy_info": version.GetBuildInfo(),
	})
}

func Run() {
	fmt.Printf("Gin start")

	r := gin.Default()

	projectPath := "/host"
	r.GET(projectPath, DefaultRoot)

	API := r.Group(projectPath)
	API.GET("/", DefaultRoot)

	defaultPort := ":80"
	r.Run(defaultPort)
}