package hub

import (
	"bytes"
	"testing"
)

func TestMain(m *testing.M) {

	m.Run()
}

func TestConnect(t *testing.T) {

	opts := HubOpts{
		Server: "127.0.0.1:4222",
		Auth:   false,
	}

	conn, err := Connect(opts)
	if err != nil {
		t.Error(err)
	}

	t.Log(conn.connName)
}

func TestOnRequested(t *testing.T) {

	opts := HubOpts{
		Server: "127.0.0.1:4222",
		Auth:   false,
	}

	conn, err := Connect(opts)
	if err != nil {
		t.Error(err)
	}

	conn.OnRequested(
		"testing",
		func(req HubReq, subj string, msg []byte) {
			if subj != "testing" {
				t.Error("invalid_subject_on_requested")
			}

			if !bytes.Equal(msg, []byte("testing")) {
				t.Error("invalid_msg_on_requested")
				t.Log(msg)
			}
		},
	)

	conn.Query("testing", []byte("testing"), 1)
}
