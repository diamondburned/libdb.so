package vars

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/internal/vars"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(app))
}

var app = cli.App{
	Name:      "vars",
	Usage:     "list, get, and set flag variables",
	UsageText: "vars [options] [command]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "always print in JSON format",
		},
	},
	Commands: []*cli.Command{
		{
			Name:      "get",
			Usage:     "get the value of a variable",
			UsageText: "vars get [options] <name>",
			Action:    get,
		},
		{
			Name:      "set",
			Usage:     "set the value of a variable",
			UsageText: "vars set [options] <name> <value>",
			Action:    set,
		},
	},
	Action: list,
}

func list(c *cli.Context) error {
	env := vm.EnvironmentFromContext(c.Context)

	variables := vars.Variables()
	if c.Bool("json") {
		b, _ := json.MarshalIndent(variables, "", "  ")
		env.Println(string(b))
		return nil
	}

	w := tabwriter.NewWriter(env.Terminal.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "[key]", "[value]", "[description]")
	for _, variable := range variables {
		var value json.RawMessage
		variable.Get(&value)
		fmt.Fprintf(w, "%s\t%s\t%s\n", variable.Key, value, variable.Description)
	}
	w.Flush()
	return nil
}

func get(c *cli.Context) error {
	env := vm.EnvironmentFromContext(c.Context)
	log := vm.LoggerFromContext(c.Context)

	name := c.Args().First()
	if name == "" {
		return cli.Exit("missing variable name", 1)
	}

	v := vars.Get(name)
	if v == nil {
		return cli.Exit(fmt.Sprintf("unknown variable %s", name), 1)
	}

	var value json.RawMessage
	ok, err := v.Get(&value)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}
	if !ok {
		log.Printf("variable %s is not set", name)
		return nil
	}

	fmt.Fprintln(env.Terminal.Stdout, string(value))
	return nil
}

func set(c *cli.Context) error {
	name := c.Args().Get(0)
	if name == "" {
		return cli.Exit("missing variable name", 1)
	}

	v := vars.Get(name)
	if v == nil {
		return cli.Exit(fmt.Sprintf("unknown variable %s", name), 1)
	}

	value := c.Args().Get(1)
	if value == "" {
		return cli.Exit("missing variable value", 1)
	}

	var jsonValue any
	if err := json.Unmarshal([]byte(value), &jsonValue); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if err := v.Set(jsonValue); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}
