package main

import (
	"flag"
	"log"
	"strings"

	"github.com/pd0mz/go-aprs"
	"github.com/pd0mz/go-aprs/aprsis"
)

func main() {
	addressString := flag.String("address", "", "calling address")
	server := flag.String("server", "euro.aprs2.net:14580", "server address")
	filter := flag.String("filter", "", "APRS filter")
	flag.Parse()

	address := aprs.MustParseAddress(*addressString)
	client, err := aprsis.Dial("tcp", *server)
	if err != nil {
		log.Fatalln(err)
	}
	if err = client.Login(address, *filter); err != nil {
		log.Fatalln(err)
	}

	packets := make(chan aprs.Packet)
	go func() {
		for {
			packet := <-packets
			if strings.Contains(strings.ToLower(packet.Raw), "dmr") {
				log.Println(packet)
			}
		}
	}()

	log.Fatalln(client.ReadPackets(packets))
}
