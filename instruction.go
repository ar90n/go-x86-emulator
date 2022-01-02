package main

import "fmt"

var instructions = make([]func(emu *Emulator), 256)

func mov_r32_imm32(emu *Emulator) {
	reg := get_code8(emu, 0) - 0xB8
	value := get_code32(emu, 1)
	set_register32(emu, Register(reg), value)
	emu.eip += 5
}

func mov_rm32_r32(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	r32 := get_r32(emu, &modrm)
	set_rm32(emu, &modrm, r32)
}

func mov_r32_rm32(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	rm32 := get_rm32(emu, &modrm)
	set_r32(emu, &modrm, rm32)
}

func mov_rm32_imm32(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	imm32 := get_code32(emu, 0)
	emu.eip += 4
	set_rm32(emu, &modrm, imm32)
}

func short_jump(emu *Emulator) {
	diff := int32(get_sign_code8(emu, 1))
	emu.eip += uint32(diff + 2)
}

func near_jump(emu *Emulator) {
	diff := get_sign_code32(emu, 1)
	emu.eip += uint32(diff + 5)
}

func init_instructions() {
	instructions[0x01] = add_rm32_r32
	instructions[0x3B] = cmp_r32_rm32
	instructions[0x3C] = cmp_al_imm8
	instructions[0x3D] = cmp_eax_imm32

	for i := 0; i < 8; i++ {
		instructions[0x40+i] = inc_r32
	}

	for i := 0; i < 8; i++ {
		instructions[0x48+i] = dec_r32
	}

	for i := 0; i < 8; i++ {
		instructions[0x50+i] = push_r32
	}

	for i := 0; i < 8; i++ {
		instructions[0x58+i] = pop_r32
	}
	instructions[0x68] = push_imm32
	instructions[0x6A] = push_imm8

	instructions[0x70] = jo
	instructions[0x71] = jno
	instructions[0x72] = jc
	instructions[0x73] = jnc
	instructions[0x74] = jz
	instructions[0x75] = jnz
	instructions[0x78] = js
	instructions[0x79] = jns
	instructions[0x7C] = jl
	instructions[0x7E] = jle

	instructions[0x83] = code_83
	instructions[0x88] = mov_rm8_r8
	instructions[0x89] = mov_rm32_r32
	instructions[0x8A] = mov_r8_rm8
	instructions[0x8B] = mov_r32_rm32
	for i := 0; i < 8; i++ {
		instructions[0xB0+i] = mov_r8_imm8
	}
	for i := 0; i < 8; i++ {
		instructions[0xb8+i] = mov_r32_imm32
	}
	instructions[0xC3] = ret
	instructions[0xC7] = mov_rm32_imm32
	instructions[0xC9] = leave

	instructions[0xCD] = swi

	instructions[0xE8] = call_rel32
	instructions[0xe9] = near_jump
	instructions[0xeb] = short_jump
	instructions[0xEC] = in_al_dx
	instructions[0xEE] = out_dx_al
	instructions[0xFF] = code_ff
}

func add_rm32_r32(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	r32 := get_r32(emu, &modrm)
	rm32 := get_rm32(emu, &modrm)
	set_rm32(emu, &modrm, rm32+r32)
}

func add_rm32_imm8(emu *Emulator, modrm *ModRM) {
	rm32 := get_rm32(emu, modrm)
	imm8 := uint32(int32(get_sign_code8(emu, 0)))
	emu.eip += 1
	set_rm32(emu, modrm, rm32+imm8)
}

func sub_rm32_imm8(emu *Emulator, modrm *ModRM) {
	rm32 := get_rm32(emu, modrm)
	imm8 := uint32(get_sign_code8(emu, 0))
	emu.eip += 1
	result := uint64(rm32) - uint64(imm8)
	set_rm32(emu, modrm, uint32(result))
	update_eflags_sub(emu, rm32, imm8, result)
}

func code_83(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)

	switch modrm.opcode_or_reg_index {
	case 0:
		add_rm32_imm8(emu, &modrm)
	case 5:
		sub_rm32_imm8(emu, &modrm)
	case 7:
		cmp_rm32_imm8(emu, &modrm)
	default:
		panic(fmt.Sprintf("not implemented: 83 %d\n", modrm.opcode_or_reg_index))
	}
}

func inc_rm32(emu *Emulator, modrm *ModRM) {
	value := get_rm32(emu, modrm)
	set_rm32(emu, modrm, value+1)
}

func code_ff(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)

	switch modrm.opcode_or_reg_index {
	case 0:
		inc_rm32(emu, &modrm)
	default:
		panic(fmt.Sprintf("not implemented: ff %d\n", modrm.opcode_or_reg_index))
	}
}

func push_r32(emu *Emulator) {
	reg := get_code8(emu, 0) - 0x50
	push32(emu, get_register32(emu, Register(reg)))
	emu.eip += 1
}

func pop_r32(emu *Emulator) {
	reg := get_code8(emu, 0) - 0x58
	set_register32(emu, Register(reg), pop32(emu))
	emu.eip += 1
}

func call_rel32(emu *Emulator) {
	diff := get_sign_code32(emu, 1)
	push32(emu, emu.eip+5)
	emu.eip += uint32(diff + 5)
}

func ret(emu *Emulator) {
	emu.eip = pop32(emu)
}

func leave(emu *Emulator) {
	ebp := get_register32(emu, EBP)
	set_register32(emu, ESP, ebp)
	set_register32(emu, EBP, pop32(emu))
	emu.eip += 1
}

func push_imm32(emu *Emulator) {
	value := get_code32(emu, 1)
	push32(emu, value)
	emu.eip += 5
}

func push_imm8(emu *Emulator) {
	value := get_code8(emu, 1)
	push32(emu, uint32(value))
	emu.eip += 2
}

func cmp_r32_rm32(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	r32 := get_r32(emu, &modrm)
	rm32 := get_rm32(emu, &modrm)
	result := uint64(r32) - uint64(rm32)
	update_eflags_sub(emu, r32, rm32, result)
}

func cmp_rm32_imm8(emu *Emulator, modrm *ModRM) {
	rm32 := get_rm32(emu, modrm)
	imm8 := uint32(get_sign_code8(emu, 0))
	emu.eip += 1
	result := uint64(rm32) - uint64(imm8)
	update_eflags_sub(emu, rm32, imm8, result)
}

func jx(emu *Emulator, pred func(emu *Emulator) bool) {
	var diff int8
	if pred(emu) {
		diff = get_sign_code8(emu, 1)
	} else {
		diff = 0
	}
	emu.eip += uint32(int32(diff + 2))
}

func jnx(emu *Emulator, pred func(emu *Emulator) bool) {
	jx(emu, func(emu *Emulator) bool { return !pred(emu) })
}

func js(emu *Emulator) {
	jx(emu, is_sign)
}

func jns(emu *Emulator) {
	jnx(emu, is_sign)
}

func jc(emu *Emulator) {
	jx(emu, is_carry)
}

func jnc(emu *Emulator) {
	jnx(emu, is_carry)
}

func jz(emu *Emulator) {
	jx(emu, is_zero)
}

func jnz(emu *Emulator) {
	jnx(emu, is_zero)
}

func jo(emu *Emulator) {
	jx(emu, is_overflow)
}

func jno(emu *Emulator) {
	jnx(emu, is_overflow)
}

func jl(emu *Emulator) {
	pred := func(emu *Emulator) bool {
		return is_sign(emu) != is_overflow(emu)
	}
	jx(emu, pred)
}

func jle(emu *Emulator) {
	pred := func(emu *Emulator) bool {
		return is_zero(emu) || (is_sign(emu) != is_overflow(emu))
	}
	jx(emu, pred)
}

func jge(emu *Emulator) {
	pred := func(emu *Emulator) bool {
		return is_sign(emu) != is_overflow(emu)
	}
	jnx(emu, pred)
}

func jg(emu *Emulator) {
	pred := func(emu *Emulator) bool {
		return is_zero(emu) || (is_sign(emu) != is_overflow(emu))
	}
	jnx(emu, pred)
}

func dec_r32(emu *Emulator) {
	reg := get_code8(emu, 0) - 0x48
	r32 := get_register32(emu, Register(reg))
	result := uint64(r32) - uint64(1)
	set_register32(emu, Register(reg), uint32(result))
	update_eflags_sub(emu, r32, 1, result)
	emu.eip += 1
}

func in_al_dx(emu *Emulator) {
	address := uint16(get_register32(emu, EDX) & 0xffff)
	value := io_in8(address)
	set_register8(emu, AL, value)
	emu.eip += 1
}

func out_dx_al(emu *Emulator) {
	address := uint16(get_register32(emu, EDX) & 0xffff)
	value := get_register8(emu, AL)
	io_out8(address, value)
	emu.eip += 1
}

func mov_r8_imm8(emu *Emulator) {
	reg := get_code8(emu, 0) - 0xb0
	set_register8(emu, Register(reg), get_code8(emu, 1))
	emu.eip += 2
}

func mov_rm8_r8(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	r8 := get_r8(emu, &modrm)
	set_rm8(emu, &modrm, r8)
}

func cmp_al_imm8(emu *Emulator) {
	value := get_code8(emu, 1)
	al := get_register8(emu, AL)
	result := uint64(al) - uint64(value)
	update_eflags_sub(emu, uint32(al), uint32(value), result)
	emu.eip += 2
}

func cmp_eax_imm32(emu *Emulator) {
	value := get_code32(emu, 1)
	eax := get_register32(emu, EAX)
	result := uint64(eax) - uint64(value)
	update_eflags_sub(emu, eax, value, result)
	emu.eip += 5
}

func inc_r32(emu *Emulator) {
	reg := get_code8(emu, 0) - 0x40
	set_register32(emu, Register(reg), get_register32(emu, Register(reg))+1)
	emu.eip += 1
}

func mov_r8_rm8(emu *Emulator) {
	emu.eip += 1
	modrm := parse_modrm(emu)
	rm8 := get_rm8(emu, &modrm)
	set_r8(emu, &modrm, rm8)
}

func swi(emu *Emulator) {
	int_index := get_code8(emu, 1)
	emu.eip += 2

	switch int_index {
	case 0x10:
		bios_video(emu)
	default:
		panic(fmt.Sprintf("unknown interrupt: 0x%02x\n", int_index))
	}
}
