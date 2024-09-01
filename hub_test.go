package hub

import (
	"testing"
)

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
