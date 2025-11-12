package mqtt

import (
	"encoding/json"
	"fmt"
	"iot-home/infrastructure/bark"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Subscribe 订阅一个 MQTT Topic
func (m *MQTTClient) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected, cannot subscribe to topic %s", topic)
	}

	token := m.client.Subscribe(topic, qos, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	log.Printf("MQTT: Subscribed to topic %s with QoS %d", topic, qos)
	return nil
}

// Unsubscribe 取消订阅一个 MQTT Topic
func (m *MQTTClient) Unsubscribe(topic string) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected, cannot unsubscribe from topic %s", topic)
	}

	token := m.client.Unsubscribe(topic)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topic %s: %w", topic, token.Error())
	}

	log.Printf("MQTT: Unsubscribed from topic %s", topic)
	return nil
}

func HandleSmsMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic: %s", msg.Topic())
	var data SmsMessage
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("Error unmarshaling ArchiveBookmarkMessage: %v, payload: %s", err, msg.Payload())
		return
	}
	log.Printf("Processing SMS: sender=%s, content=%s, operator=%s, timestamp=%s", data.Sender, data.Content, data.Operator, data.Timestamp)

	// --- 构建 Bark 消息 ---

	// 1. 标题 (Title): 突出发送者
	// 使用 emoji 让标题更生动
	barkTitle := fmt.Sprintf("✉️ 新短信来自 %s", data.Sender)

	// 2. 正文 (Body): 包含短信内容、运营商和时间
	// 使用 Markdown 格式让内容更易读 (Bark 支持 Markdown)

	barkBody := fmt.Sprintf(
		"%s\n\n"+
			"发件号码: %s\n"+
			"发件时间: %s\n\n"+
			"运营商: %s\n",
		data.Content,
		data.Sender,
		data.Timestamp.Format(time.DateTime),
		data.Operator,
	)

	// 3. 发送 Bark 通知
	err := bark.SendToBark(
		barkBody,
		bark.WithTitle(barkTitle),
		bark.WithGroup("SMS"),       // 将所有短信通知归类到 "SMS" 分组
		bark.WithLevel("active"),    // 设置为 "active" 级别，确保通知及时送达并引起注意
		bark.WithCopy(data.Content), // 点击通知时，自动复制短信内容，方便粘贴
	)

	if err != nil {
		log.Printf("Error sending SMS notification to Bark: %v", err)
	} else {
		log.Printf("SMS notification sent to Bark successfully.")
	}
}
