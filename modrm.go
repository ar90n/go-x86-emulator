package main

type OpcodeOrRegIndexTag int

const (
	Opcode OpcodeOrRegIndexTag = iota
	RegIndex
)

type DispTag int

const (
	Disp8 DispTag = iota
	Disp32
)

type ModRM struct {
	mod                     uint8
	opcode_or_reg_index_tag OpcodeOrRegIndexTag
	opcode_or_reg_index     uint8
	rm                      uint8

	sib                 uint8
	disp8_or_disp32_tag DispTag
	disp8_or_disp32     uint32
}

func parse_modrm(emu *Emulator) ModRM {
	modrm := ModRM{}

	code := get_code8(emu, 0)
	modrm.mod = (code & 0xc0) >> 6
	modrm.opcode_or_reg_index = (code & 0x38) >> 3
	modrm.rm = code & 0x07

	emu.eip += 1

	if modrm.mod != 3 && modrm.rm == 4 {
		modrm.sib = get_code8(emu, 0)
		emu.eip += 1
	}

	if (modrm.mod == 0 && modrm.rm == 5) || modrm.mod == 2 {
		modrm.disp8_or_disp32 = get_code32(emu, 0)
		modrm.disp8_or_disp32_tag = Disp32
		emu.eip += 4
	} else if modrm.mod == 1 {
		modrm.disp8_or_disp32 = uint32(get_sign_code8(emu, 0))
		modrm.disp8_or_disp32_tag = Disp8
		emu.eip += 1
	}

	return modrm
}

func set_rm32(emu *Emulator, modrm *ModRM, value uint32) {
	if modrm.mod == 3 {
		set_register32(emu, Register(modrm.rm), value)
	} else {
		address := calc_memory_address(emu, modrm)
		set_memory32(emu, address, value)
	}
}

func calc_memory_address(emu *Emulator, modrm *ModRM) uint32 {
	switch modrm.mod {
	case 0:
		switch modrm.rm {
		case 4:
			panic("nod implemented ModRM mod = 0, rm = 4")
		case 5:
			return modrm.disp8_or_disp32
		default:
			return get_register32(emu, Register(modrm.rm))
		}
	case 1:
		switch modrm.rm {
		case 4:
			panic("nod implemented ModRM mod = 1, rm = 4")
		default:
			return get_register32(emu, Register(modrm.rm)) + modrm.disp8_or_disp32
		}
	case 2:
		switch modrm.rm {
		case 4:
			panic("nod implemented ModRM mod = 2, rm = 4")
		default:
			return get_register32(emu, Register(modrm.rm)) + modrm.disp8_or_disp32
		}
	default:
		panic("nod implemented ModRM mod = 3")
	}
}

func get_rm32(emu *Emulator, modrm *ModRM) uint32 {
	switch modrm.mod {
	case 3:
		return get_register32(emu, Register(modrm.rm))
	default:
		address := calc_memory_address(emu, modrm)
		return get_memory32(emu, address)
	}
}

func set_r32(emu *Emulator, modrm *ModRM, value uint32) {
	set_register32(emu, Register(modrm.opcode_or_reg_index), value)
}

func get_r32(emu *Emulator, modrm *ModRM) uint32 {
	return get_register32(emu, Register(modrm.opcode_or_reg_index))
}

func get_r8(emu *Emulator, modrm *ModRM) uint8 {
	return get_register8(emu, Register(modrm.opcode_or_reg_index))
}

func set_rm8(emu *Emulator, modrm *ModRM, value uint8) {
	if modrm.mod == 3 {
		set_register8(emu, Register(modrm.opcode_or_reg_index), value)
	} else {
		address := calc_memory_address(emu, modrm)
		set_memory8(emu, address, uint32(value))
	}
}

func get_rm8(emu *Emulator, modrm *ModRM) uint8 {
	if modrm.mod == 3 {
		return get_register8(emu, Register(modrm.rm))
	} else {
		address := calc_memory_address(emu, modrm)
		return uint8(get_memory8(emu, address))
	}
}

func set_r8(emu *Emulator, modrm *ModRM, value uint8) {
	set_register8(emu, Register(modrm.opcode_or_reg_index), value)
}
