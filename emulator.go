package main

type Register int

const (
	EAX Register = iota
	ECX
	EDX
	EBX
	ESP
	EBP
	ESI
	EDI
	REGISTER_COUNT
	AL Register = EAX
	CL Register = ECX
	DL Register = EDX
	BL Register = EBX
	AH Register = AL + 4
	CH Register = CL + 4
	DH Register = DL + 4
	BH Register = BL + 4
)

func (r Register) String() string {
	switch r {
	case EAX:
		return "EAX"
	case ECX:
		return "ECX"
	case EDX:
		return "EDX"
	case EBX:
		return "EBX"
	case ESP:
		return "ESP"
	case EBP:
		return "EBP"
	case ESI:
		return "ESI"
	case EDI:
		return "EDI"
	}

	return ""
}

type Emulator struct {
	registers [REGISTER_COUNT]uint32
	eflags    uint32
	memory    []uint8
	eip       uint32
}
