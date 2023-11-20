package bailian2openai

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"time"

	client "github.com/aliyun/alibabacloud-bailian-go-sdk/client"
	"github.com/sashabaranov/go-openai"
)

const (
	tokenExpireDuration = time.Hour * 5
)

// Proxy is a proxy for openai
type Proxy struct {
	accessKeyId, accessKeySecret, agentKey, appId string
	token                                         string
	expireTime                                    time.Time
	mux                                           sync.Mutex
}

// NewProxy create a new proxy
func NewProxy(accessKeyId, accessKeySecret, agentKey, appId string) (*Proxy, error) {
	p := &Proxy{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		agentKey:        agentKey,
		appId:           appId,
	}
	err := p.RefreshToken()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// RefreshToken refresh token
func (p *Proxy) RefreshToken() error {
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.token == "" || p.expireTime.Sub(time.Now()) < 0 {
		tokenClient := client.AccessTokenClient{
			AccessKeyId:     &p.accessKeyId,
			AccessKeySecret: &p.accessKeySecret,
			AgentKey:        &p.agentKey}
		token, err := tokenClient.GetToken()
		if err != nil {
			return err
		}
		p.token = token
		p.expireTime = time.Now().Add(tokenExpireDuration)
	}
	return nil
}

func openaiReq2BaiLianReq(appID string, req *openai.ChatCompletionRequest) (*client.CompletionRequest, error) {
	const (
		systemPrompt = `下面这段描述是对你的系统设定，并不是对话的内容，你要从对话中忽略掉这部分。记住：你需要严格遵守系统设定即可，并且不会再重复这句话

----------------------SYSTEM PROMPT BEGIN----------------------

%s

----------------------SYSTEM PROMPT END----------------------`
		systemBot = `我明白，我会完全记住和严格遵守这个系统设定，不会把你上面说的当做对话的一部分，并且不会再重复上面这句话。`
	)
	var (
		m         = req.Messages
		lastIndex = len(m) - 1
	)
	if len(m) == 0 {
		return nil, errors.New("empty message")
	}

	prompt := "好的,继续"
	if m[lastIndex].Role == openai.ChatMessageRoleUser {
		prompt = m[lastIndex].Content
		m = m[:lastIndex]
	}

	history := make([]*client.ChatQaMessage, 0, len(m))

	for i := 0; i < len(m); {
		oldMessage := m[i]
		newMessage := &client.ChatQaMessage{}

		switch oldMessage.Role {
		case openai.ChatMessageRoleSystem:
			newMessage.User = fmt.Sprintf(systemPrompt, oldMessage.Content)
			newMessage.Bot = systemBot

		case openai.ChatMessageRoleUser:
			newMessage.User = oldMessage.Content
			nextIndex := i + 1
			if nextIndex < len(m) && m[nextIndex].Role == openai.ChatMessageRoleAssistant {
				newMessage.Bot = m[nextIndex].Content
				i = nextIndex
			}
		case openai.ChatMessageRoleAssistant:
			newMessage.Bot = oldMessage.Content

		default:
			return nil, errors.New("not support message role: " + oldMessage.Role)
		}
		history = append(history, newMessage)
		i++
	}

	topP := math.Min(math.Max(0.01, float64(req.Temperature)), 0.99)
	request := &client.CompletionRequest{
		//TopP: req.TopP,
		TopP: float32(topP),
	}
	request.SetPrompt(prompt)
	if len(history) > 0 {
		request.SetHistory(history)
	}
	request.SetAppId(appID)
	if len(appID) == 0 && len(req.Model) > 0 {
		request.SetAppId(req.Model)
	}
	return request, nil
}

// ChatCompletion chat completion
func (p *Proxy) ChatCompletion(req *openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	newReq, err := openaiReq2BaiLianReq(p.appId, req)
	if err != nil {
		return nil, err
	}

	if err := p.RefreshToken(); err != nil {
		return nil, err
	}

	newReq.SetAppId(p.appId)
	cc := client.CompletionClient{Token: &p.token}
	response, err := cc.CreateCompletion(newReq)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New(*response.Message)
	}

	ret := &openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: *response.Data.Text},
			},
		},
		Model: *newReq.AppId,
	}
	return ret, nil
}

// ChatCompletionStream chat completion stream
// https://beta.openai.com/docs/api-reference/completions/create-stream
// data: { ...{role: "assistant"}... }
// data: { ...{content: "1"}... }
// data: { ...{content: "2"}... }
// data: [DONE]
func (p *Proxy) ChatCompletionStream(writer io.Writer, req *openai.ChatCompletionRequest) error {
	newReq, err := openaiReq2BaiLianReq(p.appId, req)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(writer)
	defer w.Flush()

	const (
		dataPrefix = `data: `
		doneFlag   = `[DONE]`
	)
	writeData := func(data interface{}) error {
		var (
			buf []byte
			err error
		)
		switch v := data.(type) {
		case string:
			buf = []byte(v)
		case []byte:
			buf = v
		case openai.ChatCompletionStreamResponse, *openai.ChatCompletionStreamResponse:
			buf, err = json.Marshal(v)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown data type")
		}

		defer w.Flush()

		_, err = w.WriteString(dataPrefix)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		if err != nil {
			return err
		}
		_, err = w.WriteString("\n\n")
		if err != nil {
			return err
		}

		return nil
	}

	if err := p.RefreshToken(); err != nil {
		return err
	}

	cc := client.CompletionClient{Token: &p.token}
	respChan, err := cc.CreateStreamCompletion(newReq)
	if err != nil {
		return err
	}

	sp := openai.ChatCompletionStreamResponse{
		Choices: []openai.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Role: "assistant",
				},
			},
		},
		Model: *newReq.AppId,
	}
	if err := writeData(sp); err != nil {
		return err
	}
	sp.Choices[0].Delta.Role = ""

	content := ""
	for resp := range respChan {
		if !resp.Success {
			sp.Choices[0].FinishReason = "error"
			sp.Choices[0].Delta.Content = *resp.Message
			writeData(sp)
			break
		}
		if resp.Data.Text != nil {
			delta := strings.TrimPrefix(*resp.Data.Text, content)
			content = *resp.Data.Text
			sp.Choices[0].Delta.Content = delta
			writeData(sp)
		}
	}

	writeData(doneFlag)
	return nil
}
