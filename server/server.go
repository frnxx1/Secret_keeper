package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"server/pkg/key"
	"server/pkg/storage"

	"github.com/gin-gonic/gin"
)

const MAX_TTL = 86400
const MIN_TTL = 60

func validationCheckMessage(msg string) bool {
	return len(msg) <= 10
}

func writeInternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
}

func validationTTL(ttl int) bool {
	if ttl >= MAX_TTL {
		return false
	}
	if ttl < MIN_TTL {
		return false
	}
	return true
}

func writeLimitError(c *gin.Context) {
	c.HTML(http.StatusLengthRequired, "limit.html", gin.H{})
}

func indexView(c *gin.Context) {
	c.HTML(http.StatusOK,
		"index.html",
		gin.H{"maxTTL": MAX_TTL,
			"maxMessageLenght": 1024,
			"minTTL":           MIN_TTL})
}

func saveMessageView(c *gin.Context, keyBuilder key.KeyBuilder, keeper storage.Keeper) {
	message := c.PostForm("message")

	if !validationCheckMessage(message) {
		writeLimitError(c)
		return
	}

	ttl, err := strconv.Atoi(c.PostForm("ttl"))
	if err != nil {
		ttl = 0
	}

	if !validationTTL(ttl) {
		log.Println("Bad ttl")
		writeLimitError(c)
		return
	}

	key, err := keyBuilder.Get()
	if err != nil {
		writeInternalError(c)
		return
	}
	err = keeper.Set(key, message, ttl)
	if err != nil {
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "key.html", gin.H{"key": fmt.Sprintf("http://%s/%s", c.Request.Host, key)})
}

func readMessageView(c *gin.Context, keyBuilder key.KeyBuilder, keeper storage.Keeper) {
	key := c.Param("key")
	msg, err := keeper.Get(key)
	if err != nil {
		if err.Error() == storage.NotFoundError {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
			return
		}
		writeInternalError(c)
		return
	}

	c.HTML(http.StatusOK, "message.html", gin.H{"message": msg})

}

func buildHandler(fn func(c *gin.Context, keyBuilder key.KeyBuilder, keeper storage.Keeper), keyBuilder key.KeyBuilder, keeper storage.Keeper) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, keyBuilder, keeper)
	}
}

func getRouter(keyBuilder key.KeyBuilder, keeper storage.Keeper) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles(

		"templates/index.html",
		"templates/key.html",
		"templates/message.html",
		"templates/404.html",
		"templates/500.html",
		"templates/limit.html")

	router.GET("/", indexView)
	router.POST("/", buildHandler(saveMessageView, keyBuilder, keeper))
	router.GET("/:key", buildHandler(readMessageView, keyBuilder, keeper))
	return router
}

func main() {
	keyBuilder := key.UUIDKeyBuilder{}
	keeper := storage.GetRedisKeeper()
	router := getRouter(keyBuilder, keeper)
	router.Run("localhost:8080")
}
