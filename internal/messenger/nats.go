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

func (m *Messenger) SubscribeAsync(subject string, f func(m *nats.Msg)) {
	if m.connection == nil {
		slog.Info("No NATS connection, not subscribing", "subject", subject)
		return
	}

	slog.Info("Subscribing to NATS", "subject", subject)
	m.connection.Subscribe(subject, f)
}

func (m *Messenger) CleanUp() {
	if m.connection != nil {
		slog.Info("Closing connection to NATS")
		m.connection.Close()
	}
}
