package com

import (
	"log"

	"go.bug.st/serial"
)

type Port struct {
	Name   string
	Speed  uint
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
		log.Fatal(err)
	}
	if ports == nil {
		log.Fatal("no ports found")
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

func (p Port) SendData(mode *serial.Mode, msg string) {
	port, err := serial.Open(p.Name, mode)
	if err != nil {
		log.Fatal("cannot open port")
	}
	defer port.Close()
	n, err := port.Write([]byte(msg))
	log.Printf("Written %d bytes", n)
	if err != nil {
		log.Fatal("cannot send data")
	}
}

func (p Port) RecieveData(mode *serial.Mode) (string, error) {
	port, err := serial.Open(p.Name, mode)
	if err != nil {
		log.Fatal("cannot open port")
	}
	defer port.Close()
	buf := make([]byte, 256)
	n, err := port.Read(buf)
	if err != nil {
		log.Fatal("cannot read from port")
	}
	log.Printf("Recieved %d bytes", n)
	return string(buf), nil
}
