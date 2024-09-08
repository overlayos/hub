package hub

import (
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	ReadyForConnections = 5 * time.Second
	ReconnectBufSize    = 128 * 1024 * 1024
	DefaultTimeout      = 3 * time.Second
)

type HubOpts struct {
	Server   string
	Username string
	Password string
	Group    string
	Auth     bool
}

type HubConn struct {
	qconn    *nats.Conn
	connAddr string
	connName string
	group    string
}

func Connect(opts HubOpts) (hubconn *HubConn, err error) {

	natsOpts := NatsOpts()
	if opts.Auth {
		natsOpts = append(natsOpts, nats.UserInfo(opts.Username, opts.Password))
	}

	qconn, err := nats.Connect(opts.Server, natsOpts...)
	if err != nil {
		return nil, err
	}

	return &HubConn{
		connAddr: opts.Server,
		qconn:    qconn,
		connName: qconn.ConnectedServerName(),
		group:    opts.Group,
	}, nil
}

func (s *HubConn) RTT() (rtt uint64, err error) {

	dur, err := s.qconn.RTT()
	return uint64(dur.Nanoseconds()), err
}

func (s *HubConn) Send(subj string, msg []byte) (err error) {

	return s.qconn.Publish(subj, msg)
}

func (s *HubConn) Query(subj string, msg []byte, timeout int) (resp string, err error) {

	respmsg, err := s.qconn.Request(subj, nil, time.Duration(timeout)*time.Second)
	if err != nil {
		return "", err
	}

	resp = string(respmsg.Data)

	return resp, nil
}

func (s *HubConn) OnReceived(subj string, handler func(string, []byte)) {

	s.qconn.QueueSubscribe(subj, s.group, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
	})
}

func (s *HubConn) Close() (err error) {

	s.qconn.Drain()
	s.qconn.Close()

	return nil
}

func NatsOpts() []nats.Option {

	_uuid, _ := uuid.NewUUID()
	cid := _uuid.String()

	return []nats.Option{
		nats.Name(cid),

		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
		nats.ReconnectBufSize(ReconnectBufSize),
		nats.RetryOnFailedConnect(true),

		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
		}),
	}
}
