package main

import (
	"iot-home/api"
	"iot-home/config"
	"iot-home/infrastructure/mqtt"
	"iot-home/infrastructure/util"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Logger
	e.Use(config.RequestLogger())

	// Routes
	api.IndexRoutes(e)

	// 1. 配置 MQTT 客户端
	cfg := &mqtt.Config{
		MQTTConfig: mqtt.MQTTConfig{
			Host:     util.GetEnv("MQTT_HOST"),
			Port:     util.GetIntEnv("MQTT_PORT"),
			ClientID: "go_mqtt_listener_client_" + time.Now().Format("20060102150405"),
			QoS:      1,
			Retained: false,
		},
	}

	// 2. 创建 MQTT 客户端实例
	mqttClient, err := mqtt.NewMQTTClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	defer mqttClient.Close() // 确保在程序退出时关闭连接

	// 3. 连接到 MQTT Broker
	err = mqttClient.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	log.Println("MQTT client connected successfully.")

	// 4. 订阅你感兴趣的 Topic
	err = mqttClient.Subscribe(mqtt.TopicEspSms, 1, mqtt.HandleSmsMessage)
	if err != nil {
		log.Printf("Error subscribing to %s: %v", mqtt.TopicEspSms, err)
	}

	// Start the service
	e.Logger.Fatal(e.Start(":1323"))
}
