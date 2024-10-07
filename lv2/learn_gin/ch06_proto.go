package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"untitled/proto"
)

func main() {
	router := gin.Default()
	router.GET("/morejson", moreJson)
	router.GET("protoBuf", returnProto)
	router.GET("purejson", purejson)

	router.Run(":8083")
	// protoc --go_out=. user.proto
}

func purejson(c *gin.Context) {

	c.PureJSON(http.StatusOK, gin.H{
		"html": "<b>heelp,</b>",
	})
}

func returnProto(c *gin.Context) {
	curse := []string{"python", "go", "apple"}
	user := &proto.Teacher{
		Name:   "bobby",
		Course: curse,
	}
	c.ProtoBuf(http.StatusOK, user)
}

func moreJson(c *gin.Context) {
	var msg struct {
		Name    string `json:"user"`
		Message string
		Number  int
	}
	msg.Name = "boby"
	msg.Message = "this is a json"
	msg.Number = 20

	c.JSON(http.StatusOK, msg)
}
