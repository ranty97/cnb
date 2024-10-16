package com

import (
	"encoding/binary"
	"github.com/ranty97/cnb/internal/crc"
	"github.com/ranty97/cnb/internal/utils"
	"log"
)

const (
	Flag              = '$'
	Special           = 'v'
	escapedByte       = 0x7D
	MaxPacketDataSize = 22
)

type Packet struct {
	Flag               byte
	Special            byte
	SourceAddress      byte
	DestinationAddress byte
	Data               []byte
	FSC                uint32
}

func InitializePacket(data []byte, portNumber byte) *Packet {
	return &Packet{
		Flag:               Flag,
		Special:            Special,
		SourceAddress:      portNumber,
		DestinationAddress: 0,
		Data:               data,
		FSC:                crc.CalculateCRC(data),
	}
}

func (p *Packet) SerializePacket() []byte {
	var packet []byte
	packet = append(packet, p.Flag, p.Special, p.SourceAddress, p.DestinationAddress)
	packet = append(packet, byteStuffing(p.Data)...)
	fcsBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fcsBytes, p.FSC)
	packet = append(packet, fcsBytes...)
	return packet
}

func byteStuffing(data []byte) []byte {
	var byteStuffed []byte
	for _, b := range data {
		if b == Flag || b == escapedByte {
			byteStuffed = append(byteStuffed, escapedByte)
			log.Println(byteStuffed)
		}
		byteStuffed = append(byteStuffed, b)
	}
	if len(data) != len(byteStuffed) {
		log.Println("Byte stuffing:\n", data, " -> ", byteStuffed)
	}
	return byteStuffed
}

func SplitDataIntoPackets(data []byte, portNumber byte) ([]Packet, int) {
	var packets []Packet
	packetCount := 0

	for len(data) > 0 {
		packetSize := MaxPacketDataSize
		if len(data) < MaxPacketDataSize {
			packetSize = len(data)
		}

		packet := InitializePacket(data[:packetSize], portNumber)

		utils.InvertRandomBitWithProbability(packet.Data, 0.7)

		packets = append(packets, *packet)

		packetCount++

		data = data[packetSize:]
	}

	return packets, packetCount
}

func DeserializeStream(raw []byte, processPacket func([]byte) []byte) ([][]byte, error) {
	var packets [][]byte
	var currentData []byte
	escaped := false

	for i := 0; i < len(raw); i++ {
		if escaped {
			currentData = append(currentData, raw[i])
			escaped = false
			continue
		}

		if raw[i] == escapedByte {
			escaped = true
			continue
		}

		if raw[i] == Flag {
			// Если собраны данные, проверим их целостность
			if len(currentData) > 0 {
				packetData := currentData[:len(currentData)-4]
				fscReceived := binary.BigEndian.Uint32(currentData[len(currentData)-4:])

				// Проверка CRC
				if crc.CalculateCRC(packetData) == fscReceived {
					log.Println("Packet delivered with no mismatches")
				} else {
					log.Println("Packet delivered with mismatches.\nContents of the received packet: ", string(packetData))
					packetData = crc.RestoreBit(packetData, fscReceived)
				}
				packets = append(packets, processPacket(packetData))
				currentData = nil
			}

			// Пропускаем флаг и следующие 3 байта
			if i+3 < len(raw) {
				i += 3
			}
			continue
		}

		currentData = append(currentData, raw[i])
	}

	// Обработка последнего пакета
	if len(currentData) > 0 {
		packetData := currentData[:len(currentData)-4]
		fscReceived := binary.BigEndian.Uint32(currentData[len(currentData)-4:])

		// Проверка CRC
		if crc.CalculateCRC(packetData) == fscReceived {
			log.Println("Packet delivered with no mismatches")
		} else {
			log.Println("Packet delivered with mismatches.\nContents of the received packet: ", string(packetData))
			packetData = crc.RestoreBit(packetData, fscReceived)
		}
		packets = append(packets, processPacket(packetData))
	}

	return packets, nil
}
