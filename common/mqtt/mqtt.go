package mqtt

import "context"

type MQTT interface {
	Connect(ctx context.Context) error
	Publish(ctx context.Context, topic string, message string) error
	Subscribe(ctx context.Context, topic string, callback func(message string)) error
	Disconnect()
	IsConnected() bool
}
