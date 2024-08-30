package patlite

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(host string, port int) *Client {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr, )
	if err != nil {
		// TODO
	}

	return &Client{conn}
}

func (client *Client) GetState() (*State, error) {
	if _, err := client.conn.Write(READ); err != nil {
		return nil, err
	}
	var resp []byte
	if _, err := client.conn.Read(resp); err != nil {
		return nil, err
	}
	state, err := StateFromBytes(resp)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (client *Client) SetState(state State) error {
	data := append(WRITE_HEADER, state.Bytes()...)
	if _, err := client.conn.Write(data); err != nil {
		return err
	}
	var resp []byte
	if _, err := client.conn.Read(resp); err != nil {
		return err
	}
	if len(resp) == 0 {
		return fmt.Errorf("empty response")
	}
	switch(resp[0]) {
		case ACK:
			return nil
		case NACK:
			return fmt.Errorf("patlite returned NACK")
		default:
			return fmt.Errorf("unknown response '%s'", resp)
	}
}
