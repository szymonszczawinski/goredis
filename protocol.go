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
)

type Command interface{}

type SetCommand struct {
	key, value string
}

func ParseCommand(rawMessage string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(rawMessage))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if v.Type() == resp.Array {
			valuesArray := v.Array()
			switch valuesArray[0].String() {
			case CommandSET:
				if len(valuesArray) != 3 {
					return nil, fmt.Errorf("invalid number of variables for SET command")
				}
				return SetCommand{
						key:   valuesArray[1].String(),
						value: valuesArray[2].String(),
					},
					nil
			}

		}
	}
	return nil, nil
}
