package main

import (
	"net/http"

	"github.com/byebyebruce/bailian2openai"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// RunAPIServer run http api server
func RunAPIServer(p *bailian2openai.Proxy, addr string) error {
	router := gin.Default()
	const (
		chatCompletionPath = "/chat/completions"
	)

	router.OPTIONS(chatCompletionPath, optionHandler())
	router.POST(chatCompletionPath, chatCompletionHandler(p))

	group := router.Group("/v1")
	{
		group.OPTIONS(chatCompletionPath, optionHandler())
		group.POST(chatCompletionPath, chatCompletionHandler(p))
	}

	return router.Run(addr)
}

func optionHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Set headers for CORS
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "*")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
func chatCompletionHandler(p *bailian2openai.Proxy) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := openai.ChatCompletionRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if req.Stream {
			c.Header("Content-Type", "text/event-stream")
			if err := p.ChatCompletionStream(c.Writer, &req); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
		} else {
			resp, err := p.CreateChatCompletion(&req)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		}
	}
}
