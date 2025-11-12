package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 接口定义了 MQTT 客户端可以执行的操作
type Client interface {
	Connect() error // 连接到 MQTT Broker
	Close()         // 关闭 MQTT 连接

	// 新增：订阅消息的方法
	Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error
	Unsubscribe(topic string) error
}

// --- MQTT 客户端实现 ---
type MQTTClient struct {
	client mqtt.Client
	config *MQTTConfig // 保存配置，以便在发布时使用 QoS 和 Retained
}

// NewMQTTClient 创建一个新的 MQTT 客户端实例
func NewMQTTClient(cfg *Config) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	brokerURL := fmt.Sprintf("tcp://%s:%d", cfg.MQTTConfig.Host, cfg.MQTTConfig.Port)
	opts.AddBroker(brokerURL)
	opts.SetClientID(cfg.MQTTConfig.ClientID)
	// opts.SetUsername(cfg.MQTTConfig.Username)
	// opts.SetPassword(cfg.MQTTConfig.Password)

	// 设置连接回调函数
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	// 默认的消息处理函数，如果客户端也需要订阅消息，可以设置
	// opts.SetDefaultPublishHandler(messagePubHandler)

	// 其他可选配置
	opts.SetCleanSession(true)                    // 每次连接都创建一个新的会话
	opts.SetAutoReconnect(true)                   // 自动重连
	opts.SetConnectRetryInterval(1 * time.Second) // 重连间隔

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		client: client,
		config: &cfg.MQTTConfig,
	}, nil
}

// Connect 尝试连接到 MQTT Broker
func (m *MQTTClient) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	return nil
}

// Close 关闭 MQTT 连接
func (m *MQTTClient) Close() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250) // 250ms 的优雅关闭时间
		log.Println("MQTT: Disconnected")
	}
}

// --- 内部辅助函数：发布消息 ---
// publish 负责将 Go 结构体序列化为 JSON 并发布到指定的 Topic
func (m *MQTTClient) publish(topic string, payload interface{}) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected, cannot publish to topic %s", topic)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for topic %s: %w", topic, err)
	}

	token := m.client.Publish(topic, m.config.QoS, m.config.Retained, payloadBytes)
	token.Wait() // 等待消息发布完成
	if token.Error() != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, token.Error())
	}

	log.Printf("MQTT: Published message to topic %s: %s", topic, string(payloadBytes))
	return nil
}

// --- MQTT 连接回调函数 (来自你提供的 MQTT 示例) ---
// 注意：这些回调函数通常用于客户端订阅消息或监控连接状态。
// 对于一个纯粹的发布者客户端，它们可能不是必需的，但保留它们以供参考。

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// 这个处理函数会在客户端订阅了某个 Topic 并收到消息时被调用。
	// 如果你的客户端只发布消息，不订阅，那么这个函数可能不会被触发。
	log.Printf("MQTT: Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("MQTT: Connected")
	// 可以在这里订阅 Topic，如果需要的话
	// if token := client.Subscribe("some/topic", 1, nil); token.Wait() && token.Error() != nil {
	//     log.Printf("MQTT: Failed to subscribe: %v", token.Error())
	// }
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("MQTT: Connect lost: %v", err)
}
