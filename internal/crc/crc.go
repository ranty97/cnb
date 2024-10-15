package crc

import "fmt"

const (
	poly32 = 0xEDB88320
)

func CalculateCRC(data []byte) uint32 {
	var crc uint32 = 0xFFFFFFFF

	for _, b := range data {
		crc = crc ^ uint32(b)
		for i := 0; i < 8; i++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly32
			} else {
				crc >>= 1
			}
		}
	}

	return ^crc
}

func CheckCRC(data []byte, receivedCRC uint32) bool {
	calculatedCRC := CalculateCRC(data)

	return calculatedCRC == receivedCRC
}

func RestoreBit(data []byte, originalCRC uint32) []byte {
	for byteIndex := 0; byteIndex < len(data); byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			modifiedData := make([]byte, len(data))
			copy(modifiedData, data)

			modifiedData[byteIndex] ^= 1 << bitIndex

			modifiedCRC := CalculateCRC(modifiedData)

			if modifiedCRC == originalCRC {
				fmt.Printf("Restored bit in byte %d, bit %d\n", byteIndex, bitIndex)
				return modifiedData
			}
		}
	}
	fmt.Println("Could not restore the data")

	return data
}
