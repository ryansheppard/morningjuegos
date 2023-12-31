package messenger

import (
	"log/slog"

	"github.com/nats-io/nats.go"
)

type Messenger struct {
	connection *nats.Conn
}

func New(natsURL string) *Messenger {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		slog.Error("Error connecting to NATS", "natsURL", natsURL, "error", err)
		return &Messenger{}
	}

	return &Messenger{
		connection: conn,
	}
}

func (m *Messenger) Publish(subject string, data []byte) error {
	if m.connection == nil {
		slog.Info("No NATS connection, not publishing", "subject", subject)
		return nil
	}

	return m.connection.Publish(subject, data)
}

func (m *Messenger) PublishMessage(msg Message) error {
	bytes, err := msg.AsBytes()
	if err != nil {
		slog.Error("Failed to marshal message", "message", msg, "error", err)
		return err
	} else {
		m.Publish(msg.GetKey(), bytes)
	}

	return nil
}

func (m *Messenger) SubscribeAsync(subject string, f func(m *nats.Msg)) {
	if m.connection == nil {
		slog.Info("No NATS connection, not subscribing", "subject", subject)
		return
	}

	slog.Info("Subscribing to NATS", "subject", subject)
	m.connection.QueueSubscribe(subject, "morningjuegos", f)
}

func (m *Messenger) CleanUp() {
	if m.connection != nil {
		slog.Info("Closing connection to NATS")
		m.connection.Close()
	}
}
