package com

import (
	"errors"
	"log"
)

const Flag = 'z'
const escapedByte = 'z' + 33

type Packet struct {
	Flag               byte
	SourceAddress      byte
	DestinationAddress byte
	Data               []byte
	FSC                byte
}

func InitializePacket(data []byte) *Packet {
	return &Packet{
		Flag:               Flag,
		SourceAddress:      0,
		DestinationAddress: 0,
		Data:               data,
		FSC:                1,
	}
}

func (p *Packet) SerializePacket() []byte {
	var packet []byte
	packet = append(packet, p.Flag, p.SourceAddress, p.DestinationAddress)
	packet = append(packet, byteStuffing(p.Data)...)
	packet = append(packet, p.FSC)
	return packet
}

func DeserializePacket(raw []byte) (Packet, error) {
	if raw[0] != Flag {
		return Packet{}, errors.New("incorrect flag")
	}
	data := deByteStuffing(raw[3 : len(raw)-1])
	packet := Packet{
		Flag:               raw[0],
		SourceAddress:      raw[1],
		DestinationAddress: raw[2],
		Data:               data,
		FSC:                raw[len(raw)-1],
	}

	return packet, nil
}

func byteStuffing(data []byte) []byte {
	var byteStuffed []byte
	for _, b := range data {
		if b == byte(122) || b == escapedByte {
			byteStuffed = append(byteStuffed, escapedByte)
			log.Println(byteStuffed)
			byteStuffed = append(byteStuffed, b^0x20)
		}
		byteStuffed = append(byteStuffed, b)
		log.Println(byteStuffed)
	}
	if len(data) != len(byteStuffed) {
		log.Println("Byte stuffing:\n", data, " -> ", byteStuffed)
	}
	return byteStuffed
}

func deByteStuffing(data []byte) []byte {
	var deByteStuffed []byte
	escaped := false // flag to track if the previous byte was the escape byte

	for _, b := range data {
		if escaped {
			if b == Flag || b == escapedByte {
				deByteStuffed = append(deByteStuffed, b)
			}
			escaped = false
		} else {
			if b == escapedByte {
				escaped = true
			} else {
				deByteStuffed = append(deByteStuffed, b)
			}
		}
	}
	if len(data) != len(deByteStuffed) {
		log.Println("Byte de-stuffing:\n", data, " -> ", deByteStuffed)
	}
	return deByteStuffed
}
