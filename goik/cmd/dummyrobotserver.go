// Copyright 2025 Hans Jørgen Grimstad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"
	"os"
)

const NUM_LEGS = 6
const PACKET_LENGTH = 39

const ID_OFFSET = 0
const COXA_OFFSET = 1
const FEMUR_OFFSET = 3
const TIBIA_OFFSET = 5

type leg struct {
	id    uint8
	coxa  uint16
	femur uint16
	tibia uint16
}

type packet struct {
	networkId uint8
	legs      [NUM_LEGS]leg
}

func (p *packet) Print() {
	fmt.Printf("Network: %d\n", p.networkId)
	for _, l := range p.legs {
		fmt.Printf("  Id: %02d, C:%04d, F:%04d, T:%04d\n", l.id, l.coxa, l.femur, l.tibia)
	}
	fmt.Printf("-----------------\n")
}

func (p *packet) Unmarshal(buf []byte, n int) error {
	if n != PACKET_LENGTH {
		return fmt.Errorf("malformed packet. length: %d. expected length: %d", len(buf), PACKET_LENGTH)
	}
	p.networkId = buf[0]

	for leg := 0; leg < NUM_LEGS; leg++ {
		p.legs[leg].id = uint8(leg) + 1
		p.legs[leg].coxa = uint16(buf[COXA_OFFSET+leg*6]) + uint16(buf[COXA_OFFSET+1+leg*6])<<8
		p.legs[leg].femur = uint16(buf[FEMUR_OFFSET+leg*6]) + uint16(buf[FEMUR_OFFSET+1+leg*6])<<8
		p.legs[leg].tibia = uint16(buf[TIBIA_OFFSET+leg*6]) + uint16(buf[TIBIA_OFFSET+1+leg*6])<<8
	}

	return nil
}

func main() {

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:1337")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var p packet

	buf := make([]byte, 512)

	// Read from UDP listener in endless loop
	for {
		n, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Received: %d - %+v\n", n, buf[0:n])

		err = p.Unmarshal(buf, n)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			p.Print()
		}
	}

}
