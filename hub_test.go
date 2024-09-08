package hub

import (
	"bytes"
	"testing"
	"time"
)

var (
	opts HubOpts
	conn *HubConn
)

func TestMain(m *testing.M) {

	opts = HubOpts{
		Server: "127.0.0.1:4222",
		Auth:   false,
	}

	conn, _ = Connect(opts)

	m.Run()
}

func TestConnect(t *testing.T) {

	t.Log(conn.connName)
}

func TestOnRequested(t *testing.T) {

	go conn.OnRequested(
		"testing",
		func(subj string, msg []byte) []byte {

			if subj != "testing" {
				t.Error("invalid_subject_on_requested")
			}

			if !bytes.Equal(msg, []byte("testing")) {
				t.Error("invalid_msg_on_requested")
				t.Log(msg)
			}

			return []byte("echo.testing")
		},
	)

	time.Sleep(time.Second)

	resp, err := conn.Query("testing", []byte("testing"), 1)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(resp))
}
