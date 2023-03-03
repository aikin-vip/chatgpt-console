package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/common-nighthawk/go-figure"
	gpt3 "github.com/sashabaranov/go-gpt3"
	"os"
	"strings"
	"time"
)

// 定义生成 loading 动画的函数
func showLoading(done chan bool) {
	colors := []int{31, 33, 32, 34, 35, 36}
	colorIdx := 0
	numDots := 1
	increasing := true

	for {
		select {
		case <-done:
			fmt.Print("\r")                    // 光标回到行首
			fmt.Print(strings.Repeat(" ", 50)) // 清空行
			fmt.Print("\r")                    // 光标回到行首
			return
		default:
			color := colors[colorIdx%len(colors)]
			symbol := fmt.Sprintf("\033[0;%dm %s \033[0m", color, strings.Repeat(".", numDots))

			fmt.Printf("\r\033[0;%dmChatGPT 正在努力思考中\033[0m ", color)
			fmt.Printf("%s", symbol)
			colorIdx++

			if increasing {
				numDots++
				if numDots == 4 {
					increasing = false
				}
			} else {
				numDots--
				if numDots == 1 {
					increasing = true
				}
			}

			time.Sleep(150 * time.Millisecond)
		}
	}
}

// 定义处理回答的函数
func processResponse(mdRenderer *glamour.TermRenderer, response string) {
	mdOutput, err := mdRenderer.Render(response)
	if err != nil {
		fmt.Printf("Markdown 渲染失败: %s\n", err.Error())
		return
	}
	fmt.Println(mdOutput)
}

func main() {
	// 获取 OpenAI API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("请设置 OPENAI_API_KEY 环境变量")
		return
	}

	// 初始化 Glamour 渲染器
	renderStyle := glamour.WithEnvironmentConfig()
	mdRenderer, err := glamour.NewTermRenderer(
		renderStyle,
	)
	if err != nil {
		fmt.Println("初始化 Markdown 渲染器失败")
		return
	}

	// 输出欢迎语
	myFigure := figure.NewFigure("ChatGPT", "", true)
	colors := []string{"\033[0;31m", "\033[0;33m", "\033[0;32m", "\033[0;36m", "\033[0;34m", "\033[0;35m"}
	myFigureRainbow := ""
	for i, char := range myFigure.String() {
		if char == ' ' {
			myFigureRainbow += " "
		} else {
			myFigureRainbow += colors[i%len(colors)] + string(char) + "\033[0m"
		}
	}
	fmt.Println(myFigureRainbow)

	// 创建 ChatGPT 客户端
	client := gpt3.NewClient(apiKey)
	if err != nil {
		fmt.Printf("创建客户端失败: %s\n", err.Error())
		return
	}

	messages := []gpt3.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "你是ChatGPT, OpenAI训练的大型语言模型, 请尽可能简洁地回答我的问题",
		},
	}

	// 读取用户输入并交互
	for {
		fmt.Print("\033[0;32m请输入问题:\033[0m ")
		fmt.Print("\033[0;33m")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()              // 读取第一行数据，包含换行符
		userInput := scanner.Text() // 获取第一行数据
		for scanner.Scan() {
			userInput += "\n" + scanner.Text() // 将下一行文本添加到input中
		}
		// 将所有的字符串连接为一行，并去掉连接处的空格
		fmt.Print("\033[m")
		if userInput != "" {
			done := make(chan bool)
			go func() {
				showLoading(done)
				done <- true
			}()
			messages = append(
				messages, gpt3.ChatCompletionMessage{
					Role:    "user",
					Content: userInput,
				},
			)
			// 调用 ChatGPT API 接口生成回答
			resp, err := client.CreateChatCompletion(
				context.Background(),
				gpt3.ChatCompletionRequest{
					Model:       gpt3.GPT3Dot5Turbo0301,
					Messages:    messages,
					MaxTokens:   2048,
					Temperature: 0,
					N:           1,
				},
			)
			done <- false // 停止显示 loading 动画
			if err != nil {
				fmt.Printf("\033[0;31mChatGPT 接口调用失败: %s\033[0m\n", err.Error())
				continue
			}
			// 处理回答并输出
			processResponse(mdRenderer, resp.Choices[0].Message.Content)
			messages = append(
				messages, gpt3.ChatCompletionMessage{
					Role:    "assistant",
					Content: resp.Choices[0].Message.Content,
				},
			)
		}
	}
}
