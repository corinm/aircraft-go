package messaging

import "github.com/nats-io/nats.go"

type NatsMessaging struct {
	NatsConn *nats.Conn
}

func NewNatsMessaging(natsUrl string) (*NatsMessaging, error) {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	return &NatsMessaging{
		NatsConn: nc,
	}, nil
}

func (m *NatsMessaging) Publish(subject string, msg []byte) error {
	if m.NatsConn == nil {
		return nats.ErrConnectionClosed
	}
	return m.NatsConn.Publish(subject, msg)
}

func (m *NatsMessaging) Close() {
	if m.NatsConn != nil {
		m.NatsConn.Close()
	}
}
