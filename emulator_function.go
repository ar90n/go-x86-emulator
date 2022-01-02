package main

func get_code8(emu *Emulator, index uint32) uint8 {
	return emu.memory[emu.eip+index]
}

func get_sign_code8(emu *Emulator, index uint32) int8 {
	return int8(emu.memory[emu.eip+index])
}

func get_code32(emu *Emulator, index uint32) uint32 {
	var ret uint32 = 0x0000
	for i := 0; i < 4; i++ {
		ret |= uint32(get_code8(emu, index+uint32(i))) << (i * 8)
	}
	return ret
}

func get_sign_code32(emu *Emulator, index uint32) int32 {
	return int32(get_code32(emu, index))
}

func set_memory8(emu *Emulator, address uint32, value uint32) {
	emu.memory[address] = uint8(value & 0xff)
}

func set_memory32(emu *Emulator, address uint32, value uint32) {
	for i := 0; i < 4; i++ {
		set_memory8(emu, address+uint32(i), value>>(i*8))
	}
}

func get_register32(emu *Emulator, index Register) uint32 {
	return emu.registers[index]
}

func set_register32(emu *Emulator, index Register, value uint32) {
	emu.registers[index] = value
}

func get_memory8(emu *Emulator, address uint32) uint32 {
	return uint32(emu.memory[address])
}

func get_memory32(emu *Emulator, address uint32) uint32 {
	var ret uint32
	for i := 0; i < 4; i++ {
		ret |= get_memory8(emu, address+uint32(i)) << (8 * i)
	}
	return ret
}

func push32(emu *Emulator, value uint32) {
	address := get_register32(emu, ESP) - 4
	set_register32(emu, ESP, address)
	set_memory32(emu, address, value)
}

func pop32(emu *Emulator) uint32 {
	address := get_register32(emu, ESP)
	ret := get_memory32(emu, address)
	set_register32(emu, ESP, address+4)
	return ret
}

func update_eflags_sub(emu *Emulator, v1, v2 uint32, result uint64) {
	sign1 := (v1 >> 31) != 0
	sign2 := (v2 >> 31) != 0
	signr := ((result >> 31) & 1) != 0
	carry := (result >> 32)

	set_carry(emu, carry != 0)
	set_zero(emu, result == 0)
	set_sign(emu, signr)
	set_overflow(emu, sign1 != sign2 && sign1 != signr)
}

type Flag uint32

const (
	CARRY_FLAG    Flag = 1
	ZERO_FLAG     Flag = (1 << 6)
	SIGN_FLAG     Flag = (1 << 7)
	OVERFLOW_FLAG Flag = (1 << 11)
)

func set_flag(emu *Emulator, flag Flag, value bool) {
	if value {
		emu.eflags |= uint32(flag)
	} else {
		emu.eflags &= ^uint32(flag)
	}
}

func set_carry(emu *Emulator, is_carry bool) {
	set_flag(emu, CARRY_FLAG, is_carry)
}

func set_zero(emu *Emulator, is_zero bool) {
	set_flag(emu, ZERO_FLAG, is_zero)
}

func set_sign(emu *Emulator, is_sign bool) {
	set_flag(emu, SIGN_FLAG, is_sign)
}

func set_overflow(emu *Emulator, is_overflow bool) {
	set_flag(emu, OVERFLOW_FLAG, is_overflow)
}

func is_carry(emu *Emulator) bool {
	return (emu.eflags & uint32(CARRY_FLAG)) != 0
}

func is_zero(emu *Emulator) bool {
	return (emu.eflags & uint32(ZERO_FLAG)) != 0
}

func is_sign(emu *Emulator) bool {
	return (emu.eflags & uint32(SIGN_FLAG)) != 0
}

func is_overflow(emu *Emulator) bool {
	return (emu.eflags & uint32(OVERFLOW_FLAG)) != 0
}

func get_register8(emu *Emulator, index Register) uint8 {
	if uint32(index) < 4 {
		return uint8(emu.registers[index] & 0xff)
	} else {
		return uint8((emu.registers[index-4] >> 8) & 0xff)
	}
}

func set_register8(emu *Emulator, index Register, value uint8) {
	if uint32(index) < 4 {
		r := emu.registers[index] & 0xffffff00
		emu.registers[index] = r | uint32(value)
	} else {
		r := emu.registers[index-4] & 0xffff00ff
		emu.registers[index-4] = r | (uint32(value) << 8)
	}
}
