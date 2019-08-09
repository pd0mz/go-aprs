package aprsis

import (
	"errors"
	"log"
	"net/textproto"
	"strings"

	"github.com/hb9tf/go-aprs"
)

var (
	ErrLoginFailed = errors.New(`aprsis: login failed`)
	ErrProtocol    = errors.New(`aprsis: protocol error`)
)

type APRSIS struct {
	*textproto.Conn
}

func Dial(network, addr string) (s *APRSIS, err error) {
	s = new(APRSIS)
	s.Conn, err = textproto.Dial(network, addr)
	return
}

func (s *APRSIS) Login(address *aprs.Address, filter string) (err error) {
	if address == nil {
		return aprs.ErrAddressInvalid
	}

	if filter != "" {
		filter = " filter " + filter
	}

	if err = s.PrintfLine("user %s pass %d vers go-aprs DEV%s", address, address.Secret(), filter); err != nil {
		return
	}

	var line string
	for {
		if line, err = s.ReadLine(); err != nil {
			return
		}

		line = strings.ToLower(line)
		if strings.HasPrefix(line, "# logresp ") {
			switch {
			case strings.Index(line, " verified") > 0:
				return nil
			case strings.Index(line, " unverified") > 0:
				return ErrLoginFailed
			}
		} else if strings.HasPrefix(line, "# invalid ") {
			return ErrProtocol
		} else if strings.HasPrefix(line, "# login by user not allowed") {
			return ErrLoginFailed
		}
	}
}

func (s *APRSIS) ReadPackets(packets chan aprs.Packet) (err error) {
	var line string
	for {
		if line, err = s.ReadLine(); err != nil {
			return
		}
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		var packet aprs.Packet
		if packet, err = aprs.ParsePacket(line); err != nil {
			log.Printf("error parsing packet: %v\n", err)
			if err != aprs.ErrAddressInvalid {
				continue
			}
		}
		packets <- packet
	}
}
