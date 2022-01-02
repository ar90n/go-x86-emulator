package main

import (
	"bufio"
	"os"
)

func io_in8(address uint16) uint8 {
	switch address {
	case 0x03f8:
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		return []byte(input)[0]
	default:
		return 0
	}
}

func io_out8(address uint16, value uint8) {
	switch address {
	case 0x03f8:
		writer := bufio.NewWriter(os.Stdin)
		writer.WriteByte(value)
		writer.Flush()
	}
}
