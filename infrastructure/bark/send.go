package bark

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SendToBark 函数用于发送消息到 Bark，支持更多选项
// body: 要发送的消息正文
// options: 可选的 Bark 参数，如标题、声音、分组等
func SendToBark(body string, options ...BarkOption) error {
	// 1. 配置检查
	if appConfig.BarkAPI == "" {
		log.Printf("ERROR: util_notify: 未配置 `Config.BarkAPI`")
		return fmt.Errorf("bark API is not configured")
	}
	if appConfig.BarkKey == "" {
		log.Printf("ERROR: util_notify: 未配置 `Config.BarkKey`")
		return fmt.Errorf("bark Key is not configured")
	}

	// 2. 构建 URL
	barkURL := fmt.Sprintf("%s/%s", appConfig.BarkAPI, appConfig.BarkKey)

	// 3. 构建请求体并进行 URL 编码
	data := url.Values{}
	data.Set("body", body) // 消息正文是必须的

	// 应用所有选项
	for _, opt := range options {
		opt(&data)
	}

	// 设置默认值，如果选项中没有指定
	if data.Get("title") == "" {
		data.Set("title", "新通知") // 默认标题
	}
	if data.Get("group") == "" {
		data.Set("group", "Default") // 默认分组
	}
	if data.Get("level") == "" {
		data.Set("level", "active") // 默认级别
	}

	encodedBody := data.Encode()

	// 4. 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置请求超时时间
	}

	// 5. 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", barkURL, strings.NewReader(encodedBody))
	if err != nil {
		log.Printf("ERROR: util_notify: 创建 HTTP 请求失败: %v", err)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 6. 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 7. 发送请求
	log.Printf("INFO: util_notify: POST %s with title: %s, body: %s", barkURL, data.Get("title"), data.Get("body"))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: util_notify: 发送 HTTP 请求失败: %v", err)
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close() // 确保关闭响应体，防止资源泄露

	// 8. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		responseBodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR: util_notify: Bark 服务器返回非 OK 状态: %d, 响应: %s", resp.StatusCode, string(responseBodyBytes))
		return fmt.Errorf("bark server returned non-OK status: %d", resp.StatusCode)
	}

	log.Printf("INFO: util_notify: Bark 通知发送成功。")
	return nil // 成功
}
