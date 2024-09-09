package com

import (
	"log"
	"strings"

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
		log.Print(err)
	}
	if len(ports) == 0 {
		log.Print("no avaliable ports found")
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

func SendData(portName string, mode *serial.Mode, msg string) {
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal("cannot open port")
	}
	n, err := port.Write([]byte(msg))
	log.Printf("Written %d bytes", n)
	if err != nil {
		log.Fatal("cannot send data")
	}
	err = port.Close()
	if err != nil {
		return
	}
}

func RecieveData(portName string, mode *serial.Mode) (string, error) {
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal("cannot open port")
	}
	buf := make([]byte, 256)
	var receivedBytes = 0
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Print("cannot read from port")
			return "", err
		}
		log.Printf("Recieved %d bytes", n)
		if receivedBytes == 0 {
			log.Println("\nEOF")
			break
		}
		if strings.Contains(string(buf[:receivedBytes]), "\n") {
			break
		}
	}
	err = port.Close()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
