package programs

import "libdb.so/console"

var programs []console.Program

// Register registers the given program to the global program registry.
func Register(program console.Program) {
	programs = append(programs, program)
}

// All returns all registered programs.
func All() []console.Program {
	return programs
}
