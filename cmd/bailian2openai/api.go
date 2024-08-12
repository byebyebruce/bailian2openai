package main

import (
	"context"
	"net/http"

	"github.com/byebyebruce/bailian2openai"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// 指定部分可用模型列表
// var availableModels = []string{"gpt-3.5-turbo", "text-davinci-003", "code-davinci-002"}

// RunAPIServer run http api server
// func RunAPIServer(p *bailian2openai.Proxy, addr string) error {
// 	router := gin.Default()
// 	const (
// 		chatCompletionPath = "/chat/completions"
// 		completionPath     = "/completions"
// 	)

// 	router.OPTIONS(chatCompletionPath, optionHandler())
// 	router.POST(chatCompletionPath, chatCompletionHandler(p))
// 	router.POST(completionPath, completionHandler(p))

// 	group := router.Group("/v1")
// 	{
// 		group.OPTIONS(chatCompletionPath, optionHandler())
// 		group.POST(chatCompletionPath, chatCompletionHandler(p))
// 		group.POST(completionPath, completionHandler(p))
// 	}

//		return router.Run(addr)
//	}
func RunAPIServer(p *bailian2openai.Proxy, addr string) error {
	router := gin.Default()
	const (
		chatCompletionPath = "/chat/completions"
		completionPath     = "/completions"
	)
	router.OPTIONS(chatCompletionPath, optionHandler())
	router.POST(chatCompletionPath, chatCompletionHandler(p))
	router.POST(completionPath, completionHandler(p))
	router.GET("/models", modelsHandler()) // 添加 /models 路由
	group := router.Group("/v1")
	{
		group.OPTIONS(chatCompletionPath, optionHandler())
		group.POST(chatCompletionPath, chatCompletionHandler(p))
		group.POST(completionPath, completionHandler(p))
		group.GET("/models", modelsHandler()) // 添加 /v1/models 路由
	}
	return router.Run(addr)
}

// modelsHandler 处理 /models 路由
// modelsHandler 处理 /models 路由
func modelsHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		//这一部分适用于 12行的指定部分可用模型列表
		// data := make([]gin.H, len(availableModels))
		// for i, model := range availableModels {
		// 	data[i] = gin.H{
		// 		"id":       model,
		// 		"object":   "model",
		// 		"created":  time.Now().Unix(), // 使用当前时间戳作为创建时间
		// 		"owned_by": "system",
		// 	}
		// }
		//下面的话就是阿里百炼的模型列表
		c.JSON(http.StatusOK, gin.H{
			"object": "list",
			"data": []gin.H{
				{
					"id":       "qwen-turbo",
					"object":   "model",
					"created":  1714377100,
					"owned_by": "system",
				},
				{
					"id":       "qwen-plus",
					"object":   "model",
					"created":  1714377100,
					"owned_by": "system",
				},
				{
					"id":       "qwen-max",
					"object":   "model",
					"created":  1714377100,
					"owned_by": "system",
				},
				{
					"id":       "qwen-max-0403",
					"object":   "model",
					"created":  1714377100,
					"owned_by": "system",
				},
				{
					"id":       "qwen-max-0107",
					"object":   "model",
					"created":  1714377100,
					"owned_by": "system",
				},
				{
					"id":       "qwen-max-1201",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qwen-max-longcontext",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qwen1.5-72b-chat",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qwen1.5-32b-chat",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qwen1.5-14b-chat",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qqwen1.5-7b-chat",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qwen1.5-1.8b-chat",
					"object":   "model",
					"created":  1714377101,
					"owned_by": "system",
				},
				{
					"id":       "qwen1.5-0.5b-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "codeqwen1.5-7b-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "qwen-72b-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "qwen-14b-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "qwen-7b-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "qwen-1.8b-longcontext-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "qwen-1.8b-chat",
					"object":   "model",
					"created":  1714377102,
					"owned_by": "system",
				},
				{
					"id":       "qwen1.5-110b-chat",
					"object":   "model",
					"created":  1714788848,
					"owned_by": "system",
				},
				{
					"id":       "qwen-max-0428",
					"object":   "model",
					"created":  1714978712,
					"owned_by": "system",
				},
				{
					"id":       "qwen-vl-plus",
					"object":   "model",
					"created":  1716197934,
					"owned_by": "system",
				},
				{
					"id":       "qwen-vl-max",
					"object":   "model",
					"created":  1716197964,
					"owned_by": "system",
				},
				{
					"id":       "qwen-long",
					"object":   "model",
					"created":  1716203538,
					"owned_by": "system",
				},
				{
					"id":       "qwen2-72b-instruct",
					"object":   "model",
					"created":  1717760433,
					"owned_by": "system",
				},
				{
					"id":       "qwen2-7b-instruct",
					"object":   "model",
					"created":  1717760433,
					"owned_by": "system",
				},
				{
					"id":       "qwen2-0.5b-instruct",
					"object":   "model",
					"created":  1717760433,
					"owned_by": "system",
				},
				{
					"id":       "qwen2-1.5b-instruct",
					"object":   "model",
					"created":  1717760433,
					"owned_by": "system",
				},
				{
					"id":       "qwen2-57b-a14b-instruct",
					"object":   "model",
					"created":  1717760433,
					"owned_by": "system",
				},
			},
		})
	}
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
			if err := p.ChatCompletionStream(c.Writer, req); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
		} else {
			resp, err := p.CreateChatCompletion(context.Background(), req)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		}
	}
}

func completionHandler(p *bailian2openai.Proxy) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := openai.CompletionRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if req.Stream {
			c.Header("Content-Type", "text/event-stream")
			if err := p.CompletionStream(c.Writer, req); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
		} else {
			resp, err := p.CreateCompletion(context.Background(), req)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		}
	}
}
