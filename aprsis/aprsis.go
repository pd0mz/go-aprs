package aprsis

import (
	"errors"
	"fmt"
	"log"
	"net/textproto"
	"strings"

	"github.com/hb9tf/go-aprs"
)

const (
	version = "0.1"
)

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
			return nil, fmt.Errorf("aprsis protocol error: %s", line)
		} else if strings.HasPrefix(line, "# login by user not allowed") {
			return nil, errors.New("aprsis login failed (user not allowed)")
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
			if err != aprs.ErrAddressInvalid {
				continue
			}
		}
		packets <- packet
	}
}
