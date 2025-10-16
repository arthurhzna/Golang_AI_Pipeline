package mqtt

import (
	"context"
	"fmt"
	"time"

	errWrap "task_queue/common/error"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTImpl struct {
	client   mqtt.Client
	broker   string
	port     string
	clientID string
	username string
	password string
}

func NewMQTT(broker, port, clientID, username, password string) MQTT {
	return &MQTTImpl{
		broker:   broker,
		port:     port,
		clientID: clientID,
		username: username,
		password: password,
	}
}

func (m *MQTTImpl) Connect(ctx context.Context) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", m.broker, m.port))
	opts.SetClientID(m.clientID)

	if m.username != "" {
		opts.SetUsername(m.username)
		opts.SetPassword(m.password)
	}

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	// if strings.HasPrefix(m.broker, "ssl://") || strings.HasPrefix(m.broker, "tls://") {
	// 	opts.SetTLSConfig(&tls.Config{
	// 		InsecureSkipVerify: true,
	// 	})
	// }

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("MQTT Connection lost: %v\n", err)
	}

	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("MQTT Connected successfully")
	}

	m.client = mqtt.NewClient(opts)

	token := m.client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return errWrap.WrapError(fmt.Errorf("connection timeout"))
	}

	if err := token.Error(); err != nil {
		return errWrap.WrapError(fmt.Errorf("failed to connect to MQTT broker: %w", err))
	}

	return nil
}

func (m *MQTTImpl) Publish(ctx context.Context, topic string, message string) error {
	if m.client == nil || !m.client.IsConnected() {
		return errWrap.WrapError(fmt.Errorf("MQTT client is not connected"))
	}

	token := m.client.Publish(topic, 1, true, message)
	if !token.WaitTimeout(5 * time.Second) {
		return errWrap.WrapError(fmt.Errorf("publish timeout"))
	}

	if err := token.Error(); err != nil {
		return errWrap.WrapError(fmt.Errorf("failed to publish message: %w", err))
	}

	fmt.Printf("Message published to topic '%s': %s\n", topic, message)
	return nil
}

func (m *MQTTImpl) Subscribe(ctx context.Context, topic string, callback func(message string)) error {
	if m.client == nil || !m.client.IsConnected() {
		return errWrap.WrapError(fmt.Errorf("MQTT client is not connected"))
	}

	token := m.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		callback(string(msg.Payload()))
	})

	if !token.WaitTimeout(5 * time.Second) {
		return errWrap.WrapError(fmt.Errorf("subscribe timeout"))
	}

	if err := token.Error(); err != nil {
		return errWrap.WrapError(fmt.Errorf("failed to subscribe to topic: %w", err))
	}

	fmt.Printf("Subscribed to topic '%s'\n", topic)
	return nil
}

func (m *MQTTImpl) Disconnect() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
		fmt.Println("MQTT Disconnected")
	}
}

func (m *MQTTImpl) IsConnected() bool {
	return m.client != nil && m.client.IsConnected()
}
