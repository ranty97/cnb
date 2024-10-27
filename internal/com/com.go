package com

import (
	"github.com/ranty97/cnb/internal/utils"
	"go.bug.st/serial"
	"log"
	"math/rand"
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
	p.SendBytes(packet.SerializePacket())
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
	data, err := p.ReceiveBytes()
	if err != nil {
		return [][]byte{}, err
	}
	//could use processing func (ex. destuffing)
	return DeserializeStream(data, func(bytes []byte) []byte {
		return bytes
	})
}
