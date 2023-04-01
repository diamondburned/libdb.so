package nsfw

import (
	"context"
	"errors"
	"fmt"
	"log"

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

func (prog) Run(ctx context.Context, env vm.Environment, args []string) error {
	if nsfw.IsEnabled() {
		log.Println("nsfw is already enabled")
		return nil
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
