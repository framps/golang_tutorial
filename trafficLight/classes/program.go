package classes

// lamp colors
const (
	off = iota
	green
	yellow
	red
	redyellow
)

// Phase consists of light and number of ticks to flash the lights
type Phase struct {
	Lights int
	Ticks  int
}

// Program - Has phases and a state (active phase)
type Program struct {
	Phases []Phase
	state  int
}

// TestProgram - Turn every lamp on
var TestProgram = Program{
	Phases: []Phase{
		Phase{red, 1},
		Phase{yellow, 1},
		Phase{green, 1},
		Phase{yellow, 1},
	},
}

// WarningProgram - Traffic light is not working, just blink
var WarningProgram = Program{
	Phases: []Phase{
		Phase{yellow, 1},
		Phase{off, 1},
	},
}

// NormalProgram - Just the common traffic light
var NormalProgram = Program{
	Phases: []Phase{
		Phase{green, 3},
		Phase{yellow, 1},
		Phase{red, 3},
		Phase{redyellow, 1},
	},
}
