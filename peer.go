package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	connection     net.Conn
	messageChannel chan Message
}

func NewPeer(connection net.Conn, messageChannel chan Message) *Peer {
	return &Peer{
		connection:     connection,
		messageChannel: messageChannel,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.connection.Write(msg)
}

func (p *Peer) ReadLoop() error {
	rd := resp.NewReader(p.connection)
	for {
		value, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if value.Type() == resp.Array {
			valuesArray := value.Array()
			var cmd Command
			switch valuesArray[0].String() {
			case CommandClient:
				cmd = ClientCommand{
					value: valuesArray[1].String(),
				}
			case CommandSET:
				if len(valuesArray) != 3 {
					return fmt.Errorf("invalid number of variables for SET command")
				}
				cmd = SetCommand{
					key:   valuesArray[1].Bytes(),
					value: valuesArray[2].Bytes(),
				}
				slog.Info("got SET cmd", "cmd", cmd)
			case CommandGET:
				if len(valuesArray) != 2 {
					return fmt.Errorf("invalid number of variables for GET command")
				}

				cmd = GetCommand{
					key: valuesArray[1].Bytes(),
				}
				slog.Info("got GET cmd", "cmd", cmd)
			case CommandHELLO:
				cmd = HelloCommand{
					value: valuesArray[1].String(),
				}
			default:
				fmt.Println("the value string angle => ", valuesArray[0])
				fmt.Printf("got unknown command => %v\r\n", value.Array())
			}
			p.messageChannel <- Message{
				cmd:  cmd,
				peer: p,
			}

		}
	}
	return nil
}
