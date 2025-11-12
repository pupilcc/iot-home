package bark

import (
	"iot-home/infrastructure/util"
	"net/url"
)

// Config 结构体用于存储配置信息
type Config struct {
	BarkAPI string
	BarkKey string
}

// 全局配置实例 (在实际应用中，你可能通过配置文件加载或依赖注入)
var appConfig = Config{
	BarkAPI: util.GetEnv("BARK_API"),
	BarkKey: util.GetEnv("BARK_KEY"),
}

// BarkOption 是一个函数类型，用于修改 url.Values
type BarkOption func(values *url.Values)

// WithTitle 设置 Bark 通知标题
func WithTitle(title string) BarkOption {
	return func(values *url.Values) {
		values.Set("title", title)
	}
}

// WithGroup 设置 Bark 通知分组
// 相同 group 的通知会在 Bark App 中归类
func WithGroup(group string) BarkOption {
	return func(values *url.Values) {
		values.Set("group", group)
	}
}

// WithLevel 设置 Bark 通知级别
// 例如: "active", "timeSensitive", "passive"
func WithLevel(level string) BarkOption {
	return func(values *url.Values) {
		values.Set("level", level)
	}
}

// WithCopy 设置 Bark 通知点击后复制的内容
func WithCopy(text string) BarkOption {
	return func(values *url.Values) {
		values.Set("copy", text)
	}
}
