package com

import (
	"fmt"
	"github.com/ranty97/cnb/internal/collision"
	"github.com/ranty97/cnb/internal/utils"
	"go.bug.st/serial"
	"log"
	"math/rand"
	"time"
)

const (
	jamSignal byte = 0xFF
	slotTime       = 10 * time.Millisecond
)

type Port struct {
	Name   string
	Speed  int
	Parity serial.Parity
}

var ParityMap = map[string]serial.Parity{
	"No parity":    serial.NoParity,
	"Odd parity":   serial.OddParity,
	"Even parity":  serial.EvenParity,
	"Mark parity":  serial.MarkParity,
	"Space parity": serial.SpaceParity,
}

var Speeds = []int{110, 9600, 115200}

func GetPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Print(err)
	}
	if len(ports) == 0 {
		log.Print("no available ports found")
	}

	return ports
}

func GetParities(p map[string]serial.Parity) []string {
	keys := make([]string, 0, len(p))
	for i := range p {
		keys = append(keys, i)
	}

	return keys
}

func (p Port) SendBytes(data []byte) {
	port, err := serial.Open(p.Name, &serial.Mode{BaudRate: p.Speed, Parity: p.Parity})
	if err != nil {
		log.Println("cannot open port")
	}
	n, err := port.Write(data)
	log.Println("Written", n, "bytes")

	err = port.Close()
	if err != nil {
		return
	}
}

func (p Port) SendPacket(packet Packet) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		log.Println(i)
		data := collision.RandomlyAddCollision(packet.SerializePacket())
		p.SendBytes(data)

		if len(data) != len(packet.SerializePacket()) {
			log.Println("Collision occurred")
			jamSlice := append([]byte{}, jamSignal)
			p.SendBytes(jamSlice)
			k := r.Intn(i + 1)
			log.Println("Waiting...")
			time.Sleep(time.Duration(slotTime * (1 << k)))
		} else {
			log.Println("Packet sent with no collision")
			return
		}
	}
	log.Println("Failed to send packet after 10 attempts")
}

func (p Port) SendData(data []byte) int {
	portNumber, _ := utils.LastCharacterAsNumber(p.Name)
	packets, packetCount := SplitDataIntoPackets(data, byte(portNumber))

	for {
		if rand.Intn(2) == 0 {
			for _, packet := range packets {
				p.SendPacket(packet)
				log.Println(packet)
			}
			break
		} else {
			log.Println("the transmission channel is busy")
		}
	}

	log.Printf("Sent %d packets\n", packetCount)
	return packetCount
}

func (p Port) ReceiveBytes() ([]byte, error) {
	port, err := serial.Open(p.Name, &serial.Mode{BaudRate: p.Speed, Parity: p.Parity})

	if err != nil {
		log.Fatal("cannot open port")
	}

	buff := make([]byte, 256)
	var n, readErr = port.Read(buff)
	if readErr != nil {
		return nil, err
	}

	err = port.Close()

	log.Printf("Recieved %d bytes\n", n)
	return buff[:n], nil
}

func (p Port) ReceivePacket() ([][]byte, error) {
	var packets [][]byte

	for i := 0; i < 10; i++ {
		data, err := p.ReceiveBytes()
		log.Println("data : ", data)
		if err != nil {
			return nil, err
		}
		log.Println(data)
		if data == nil {
			break
		}

		if len(data) > 0 && data[len(data)-1] == jamSignal {
			log.Println("Received jam signal, ignoring the packet...")
			continue
		}

		result, err := DeserializeStream(data, func(bytes []byte) []byte {
			return bytes
		})
		if err != nil {
			log.Println("Failed to deserialize data, trying to receive more...")
			continue
		}

		packets = append(packets, result...)
		return packets, nil
	}

	if len(packets) == 0 {
		return nil, fmt.Errorf("no packets received")
	}
	return packets, nil
}
