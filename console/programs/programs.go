package programs

import "libdb.so/console"

var programs = map[string]console.Program{}

// Register registers the given program to the global program registry.
func Register(program console.Program) {
	if program == nil {
		return
	}

	if _, ok := programs[program.Name()]; ok {
		panic("program already registered: " + program.Name())
	}

	programs[program.Name()] = program
}

// All returns all registered programs.
func All() map[string]console.Program {
	return programs
}
