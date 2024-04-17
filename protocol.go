package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
	CommandGET = "GET"
)

type Command interface{}

type SetCommand struct {
	key, value []byte
}

type GetCommand struct {
	key []byte
}

func ParseCommand(rawMessage string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(rawMessage))
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
			switch valuesArray[0].String() {
			case CommandSET:
				if len(valuesArray) != 3 {
					return nil, fmt.Errorf("invalid number of variables for SET command")
				}
				return SetCommand{
						key:   valuesArray[1].Bytes(),
						value: valuesArray[2].Bytes(),
					},
					nil
			case CommandGET:
				if len(valuesArray) != 2 {
					return nil, fmt.Errorf("invalid number of variables for GET command")
				}

				return GetCommand{
					key: valuesArray[1].Bytes(),
				}, nil
			}

		}
	}
	return nil, nil
}
