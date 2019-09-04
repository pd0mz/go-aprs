package aprsis

import (
	"errors"
	"fmt"
	"log"
	"net/textproto"
	"strings"

	"github.com/pd0mz/go-aprs"
)

const (
	version = "0.1"
)

var ErrNotAllowed = errors.New("aprsis: login failed")

type ProtocolError struct {
	Line string
}

func (err ProtocolError) Error() string {
	return fmt.Sprintf("aprsis: protocol error: %s", err.Line)
}

func Connect(proto, addr, call, filter string) (*textproto.Conn, error) {
	conn, err := textproto.Dial(proto, addr)
	if err != nil {
		return nil, err
	}

	if filter != "" {
		filter = " filter " + filter
	}

	if err := conn.PrintfLine("user %s pass -1 vers go-aprs %s%s", call, version, filter); err != nil {
		return nil, err
	}
	for {
		line, err := conn.ReadLine()
		if err != nil {
			return nil, err
		}
		line = strings.ToLower(line)
		if strings.HasPrefix(line, "# logresp ") {
			return conn, nil
		} else if strings.HasPrefix(line, "# invalid ") {
			return nil, ProtocolError{line}
		} else if strings.HasPrefix(line, "# login by user not allowed") {
			return nil, ErrNotAllowed
		}
	}
}

func ReadPackets(conn *textproto.Conn, packets chan aprs.Packet) error {
	for {
		line, err := conn.ReadLine()
		if err != nil {
			return err
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		packet, err := aprs.ParsePacket(line)
		if err != nil {
			log.Printf("error parsing packet: %v\n", err)
			continue
		}
		packets <- packet
	}
}
