package main

import (
	"fmt"
	"os"
)

func create_emu(size uint32, eip uint32, esp uint32) Emulator {
	var emu = Emulator{
		memory: make([]uint8, size),
		eip:    eip,
	}
	emu.registers[ESP] = esp
	return emu
}

func dump_registers(emu *Emulator) {
	for i := 0; i < int(REGISTER_COUNT); i++ {
		fmt.Printf("%s = %08x\n", Register(i).String(), emu.registers[Register(i)])
	}
	fmt.Printf("EIP = %08x\n", emu.eip)
}

func parse_args(args []string) ([]string, bool) {
	ret := make([]string, 0)
	quiet := false
	for i := 0; i < len(args); i += 1 {
		if args[i] == "-q" {
			quiet = true
		} else {
			ret = append(ret, args[i])
		}
	}

	return ret, quiet
}

func main() {
	args, quiet := parse_args(os.Args)
	if len(args) != 2 {
		fmt.Println("usage: px86 filename")
		return
	}

	const MEMORY_SIZE = 1024 * 1024
	var emu = create_emu(MEMORY_SIZE, 0x7c00, 0x7c00)
	fp, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	fp.Read(emu.memory[0x7c00:])

	init_instructions()

	for emu.eip < MEMORY_SIZE {
		code := get_code8(&emu, 0)
		if !quiet {
			fmt.Printf("EIP = %X, Code = %02X\n", emu.eip, code)
		}

		if instructions[code] == nil {
			fmt.Printf("\n\nNot Implemented: %02X\n", code)
			break
		}

		instructions[code](&emu)

		if emu.eip == 0x0000 {
			fmt.Println("\n\nend of program.\n")
			break
		}
	}

	dump_registers(&emu)
}
