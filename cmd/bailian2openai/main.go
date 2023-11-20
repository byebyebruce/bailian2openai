package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/byebyebruce/bailian2openai"
	"github.com/joho/godotenv"
)

var (
	Version string = "unknown"
)
var (
	addr    = flag.String("address", ":8080", "web server listen address")
	prompt  = flag.String("prompt", "你是一个友善的人工智能助手，你会尽力解答我的任何问题", "system prompt")
	cliMode = flag.Bool("cli", false, "CLI mode")
)

func main() {
	flag.Parse()
	fmt.Println("Version:", Version)

	godotenv.Overload()

	accessKeyId := mustEnv("ACCESS_KEY_ID")
	accessKeySecret := mustEnv("ACCESS_KEY_SECRET")
	agentKey := mustEnv("AGENT_KEY")
	appId := strings.TrimSpace(os.Getenv("APP_ID")) // 可以后面请求时再传入，用的model这个参数

	p, err := bailian2openai.NewProxy(accessKeyId, accessKeySecret, agentKey, appId)
	if err != nil {
		log.Fatalln(err)
	}

	if *cliMode {
		CLI(p, *prompt)
		return
	}
	if err := RunAPIServer(p, *addr); err != nil {
		log.Fatalln(err)
	}
}

func mustEnv(env string) string {
	val := strings.TrimSpace(os.Getenv(env))
	if len(val) == 0 {
		log.Fatalf("env %s is empty", env)
	}
	return val
}
