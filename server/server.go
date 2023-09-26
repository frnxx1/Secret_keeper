package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	return ttl <= MAX_TTL && ttl > MIN_TTL
}

func writeLimitError(c *gin.Context) {
	c.HTML(http.StatusLengthRequired, "limit.html", gin.H{})
}

func indexView(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"maxTTL":           MAX_TTL,
			"maxMessageLenght": 1024,
			"minTTL":           MIN_TTL,
		},
	)
}

func saveMessageView(ctx *gin.Context, keyBuilder key.KeyBuilder, keeper keeper) {
	message := ctx.PostForm("message")

	if !validationCheckMessage(message) {
		writeLimitError(ctx)
		return
	}

	ttl, err := strconv.Atoi(ctx.PostForm("ttl"))
	if err != nil {
		ttl = 0
	}

	if !validationTTL(ttl) {
		log.Println("Bad ttl")
		writeLimitError(ctx)
		return
	}

	k, err := keyBuilder.Get()
	if err != nil {
		writeInternalError(ctx)
		return
	}

	err = keeper.Set(ctx, k, message, ttl)
	if err != nil {
		writeInternalError(ctx)
		return
	}

	ctx.HTML(http.StatusOK, "key.html", gin.H{"key": fmt.Sprintf("http://%s/%s", ctx.Request.Host, k)})
}

func readMessageView(ctx *gin.Context, _ key.KeyBuilder, keeper keeper) {
	k := ctx.Param("key")
	msg, err := keeper.Get(ctx, k)
	if err != nil {
		if err.Error() == storage.NotFoundError {
			ctx.HTML(http.StatusNotFound, "404.html", gin.H{})
			return
		}

		writeInternalError(ctx)
		return
	}

	ctx.HTML(http.StatusOK, "message.html", gin.H{"message": msg})
}

func buildHandler(fn func(ctx *gin.Context, keyBuilder key.KeyBuilder, keeper keeper), keyBuilder key.KeyBuilder, keeper keeper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fn(ctx, keyBuilder, keeper)
	}
}

type keeper interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, message string, ttl int) error
}

func getRouter(keyBuilder key.KeyBuilder, keeper keeper) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/key.html",
		"templates/message.html",
		"templates/404.html",
		"templates/500.html",
		"templates/limit.html",
	)

	router.GET("/", indexView)
	router.POST("/", buildHandler(saveMessageView, keyBuilder, keeper))
	router.GET("/:key", buildHandler(readMessageView, keyBuilder, keeper))

	return router
}

func main() {
	keyBuilder := key.NewUUIDKeyBuilder()
	keeper := storage.GetRedisKeeper()
	router := getRouter(keyBuilder, keeper)
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}
}
