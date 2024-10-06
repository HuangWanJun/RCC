package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()
	goodsGroup := router.Group("/goods")
	{
		goodsGroup.GET("/", goodsList)
		goodsGroup.GET("/:id/:action", goodsDetail)
		goodsGroup.POST("", createGoods)
	}

	v1 := router.Group("/v1")
	{
		v1.POST("login", loginEndPoint)

	}

	v2 := router.Group("/v2")
	{
		v2.POST("login", loginEndPoint)

	}

	router.Run(":8083")

}

func loginEndPoint(context *gin.Context) {

}

func createGoods(context *gin.Context) {

}

func goodsDetail(context *gin.Context) {

	id := context.Param("id")
	action := context.Param("action")

	context.JSON(http.StatusOK, gin.H{
		"id":     id,
		"action": action,
	})
}

func goodsList(context *gin.Context) {

	context.JSON(http.StatusOK, gin.H{
		"name": "goodslist",
	})
}
