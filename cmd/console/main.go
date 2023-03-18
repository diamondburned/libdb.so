package main

import (
	"context"
	"log"
	"os"

	"github.com/fatih/color"
	"libdb.so/console"
	"libdb.so/console/programs"

	_ "libdb.so/console/programs/hewwo"
)

var colors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgWhite,
}

func main() {
	io := console.IO{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	inst, err := console.NewInterpreter(io, console.InterpreterOpts{
		MakeRaw:  false, // maybe?
		Programs: programs.All(),
	})
	if err != nil {
		log.Fatalln("failed to create console instance:", err)
	}

	if err := inst.Run(context.Background()); err != nil {
		log.Fatalln("failed to run console instance:", err)
	}
}
