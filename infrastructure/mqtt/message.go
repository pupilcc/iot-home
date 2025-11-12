package mqtt

import "time"

type SmsMessage struct {
	Sender    string    `json:"sender"`    // 发送者电话号码
	Content   string    `json:"content"`   // 消息内容
	Operator  string    `json:"operator"`  // 运营商
	Timestamp time.Time `json:"timestamp"` // 时间戳，使用 time.Time 类型自动解析 ISO 8601 格式
}

type DeviceMessage struct {
	Message string `json:"message"`
}
