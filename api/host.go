package api

/*
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ListHostedZone(c *gin.Context) {
	var req vo.ListHostedZoneReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		fmt.Println("list lost", err)

		return
	}

	resp := vo.ListHostedZone(&req)

	c.JSON(http.StatusOK, gin.H{
		"data": resp})
}

func CreateHostedZone(c *gin.Context) {
	var req vo.CreateHostedZoneReq

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, err)
		fmt.Println("xxxxbad")

		return
	}

	resp := vo.CreateHostedZone(&req)

	c.JSON(http.StatusOK, gin.H{
		"data": resp})
}
 */