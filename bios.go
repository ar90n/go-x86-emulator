package main

import "fmt"

func put_string(s string) {
	for _, c := range s {
		io_out8(0x03f8, uint8(c))
	}
}

func bios_to_termianl(color uint8) uint8 {
	return [8]uint8{30, 34, 32, 36, 31, 35, 33, 37}[color&0x07]
}

func bios_video_teletype(emu *Emulator) {
	color := get_register8(emu, BL) & 0x0f
	ch := get_register8(emu, AL)
	terminal_color := bios_to_termianl(color)
	bright := (color & 0x08) >> 3
	put_string(fmt.Sprintf("\x1b[%d;%dm%c\x1b[0m", bright, terminal_color, ch))
}

func bios_video(emu *Emulator) {
	f := get_register8(emu, AH)
	switch f {
	case 0x0e:
		bios_video_teletype(emu)
	default:
		panic(fmt.Sprintf("not implemented BIOS video function 0x%02x\n", f))
	}
}
