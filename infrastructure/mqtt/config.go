package mqtt

const (
	TopicEspSms    = "esp32/sms"
	TopicEspDevice = "esp32/device"
)

// MQTTConfig 包含 MQTT 连接所需的配置
type MQTTConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	ClientID string // MQTT 客户端ID，必须唯一
	QoS      byte   // 默认的 QoS (Quality of Service) 等级，0, 1, 2
	Retained bool   // 默认的 Retained 标志
}

// Config 包含所有队列相关的配置，现在只包含 MQTTConfig
type Config struct {
	MQTTConfig MQTTConfig
	// 如果有其他配置，可以继续添加
}
