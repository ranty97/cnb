package com

import (
	"log"
)

const (
	Flag              = '$'
	Special           = 'v'
	escapedByte       = 0x7D
	MaxPacketDataSize = 22
	FCS               = 1
)

type Packet struct {
	Flag               byte
	Special            byte
	SourceAddress      byte
	DestinationAddress byte
	Data               []byte
	FSC                byte
}

func InitializePacket(data []byte, portNumber byte) *Packet {
	return &Packet{
		Flag:               Flag,
		Special:            Special,
		SourceAddress:      portNumber,
		DestinationAddress: 0,
		Data:               data,
		FSC:                FCS,
	}
}

func (p *Packet) SerializePacket() []byte {
	var packet []byte
	packet = append(packet, p.Flag, p.Special, p.SourceAddress, p.DestinationAddress)
	packet = append(packet, byteStuffing(p.Data)...)
	packet = append(packet, p.FSC)
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
			// Если мы видим байт после escaped, добавляем его к текущим данным
			currentData = append(currentData, raw[i])
			escaped = false // Сбрасываем флаг escaped
			continue
		}

		if raw[i] == escapedByte {
			// Устанавливаем флаг, что следующий байт будет застаффленным
			escaped = true
			currentData = append(currentData, raw[i]) // Добавляем escapedByte в текущие данные
			continue
		}

		if raw[i] == Flag {
			// Если текущие данные не пустые, добавляем их как пакет, пропуская байт 1
			if len(currentData) > 0 {
				// Пропускаем последний байт, если он равен 1 (FCS)
				if currentData[len(currentData)-1] == FCS {
					currentData = currentData[:len(currentData)-1] // Убираем байт 1
				}
				// Обрабатываем пакет
				packets = append(packets, processPacket(currentData))
				currentData = nil // Сбрасываем текущие данные для нового пакета
			}
			// Пропускаем флаг и следующие три байта
			i += 3
			continue
		}

		// Добавляем байт в текущие данные
		currentData = append(currentData, raw[i])
	}

	// Проверяем, остались ли данные после последнего флага
	if len(currentData) > 0 {
		if currentData[len(currentData)-1] == FCS {
			currentData = currentData[:len(currentData)-1] // Убираем байт 1
		}
		// Обрабатываем последний пакет
		packets = append(packets, processPacket(currentData))
	}

	return packets, nil
}
