package programs

import "libdb.so/vm"

var programs = map[string]vm.Program{}

// Register registers the given program to the global program registry.
func Register(program vm.Program) {
	if program == nil {
		return
	}

	if _, ok := programs[program.Name()]; ok {
		panic("program already registered: " + program.Name())
	}

	programs[program.Name()] = program
}

// All returns all registered programs.
func All() map[string]vm.Program {
	return programs
}
