package nsfw

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	_ "embed"

	"libdb.so/internal/nsfw"
	"libdb.so/vm"
	"libdb.so/vm/programs"
)

//go:embed prompt.txt
var prompt string

func init() {
	programs.Register(prog{})
}

type prog struct{}

func (prog) Name() string { return "nsfw" }

func (prog) Usage() string {
	return "opt in to potentially nsfw content"
}

var errAlreadyEnabled = errors.New("nsfw is already enabled")

func (prog) Run(ctx context.Context, env vm.Environment, args []string) error {
	if len(args) > 1 {
		switch args[1] {
		case "enable":
			break
		case "disable":
			nsfw.Disable()
			return nil
		case "get":
			env.Println(strconv.FormatBool(nsfw.IsEnabled()))
			return nil
		case "help", "-h", "--help":
			return &vm.UsageError{Usage: "nsfw [enable|get|help]"}
		default:
			return fmt.Errorf("unknown subcommand %q", args[1])
		}
	}

	if nsfw.IsEnabled() {
		return vm.WrapError(1, errAlreadyEnabled)
	}

	fmt.Fprint(env.Terminal.Stdout, prompt)

	answer, err := env.PromptLine("To agree, type `yes' or `y': ")
	if err != nil {
		return err
	}

	switch answer {
	case "yes", "y":
		nsfw.Enable()
		fmt.Println("yay! owo")
		return nil
	default:
		return errors.New("user declined")
	}
}
