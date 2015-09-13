package cobra

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/pflag"
)

type zshCmd struct {
	Name           string
	Short          string
	Formatted      string
	HasSubCommands bool
	SubCommands    []*zshCmd
	FmtFlags       []string
}

const zshCompTemp = `#compdef {{ .Name }}
_{{ .Name }}(){
	local ret=1
	_arguments -C \
		'1: :_{{ .Name }}_cmds' \
		'*::arg:->args' \
	&& ret=0
	case $state in
		(args)
			case $line[1] in
				{{ range .SubCommands }}({{ .Name }})
					{{if .HasSubCommands}}_arguments -C \
						{{ range .SubCommands }}{{ .Formatted }}{{ end }}
					&& ret=0{{ else }}_message 'no additional commands'{{ end }}
				;;
				{{ end }}(*)
					_message 'Unknown subcommand'
				;;
			esac
		esac
	return 0
}

_{{ .Name }}_cmds(){
	local commands; commands=(
		{{ range .SubCommands }}'{{ .Name }}:{{ .Short }}'
		{{ end }})
	_describe -t commands '{{ .Name }} commands' commands "$@"
}
_{{.Name}}`

func parseCommand(c *Command) *zshCmd {
	root := &zshCmd{
		Name:           c.Name(),
		Short:          c.Short,
		Formatted:      fmt.Sprintf("%s:%s\n", c.Name(), c.Short),
		HasSubCommands: false,
	}

	if len(c.Commands()) != 0 {
		root.HasSubCommands = true
		for _, cmd := range c.Commands() {
			root.SubCommands = append(root.SubCommands, parseCommand(cmd))
		}
	}
	c.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) {
		str := fmt.Sprintf("--(%s)--%s=[%s] \\\n", flag.Shorthand, flag.Name, flag.Usage)
		root.FmtFlags = append(root.FmtFlags, str)
	})
	return root
}

func (c *Command) GenZshCompletion(out *bytes.Buffer) error {
	t, err := template.New("zshComp").Parse(zshCompTemp)
	if err != nil {
		return err
	}
	cmd := parseCommand(c)
	return t.Execute(out, cmd)
}

func (c *Command) GenZshCompletionFile(filename string) error {
	out := new(bytes.Buffer)
	err := c.GenZshCompletion(out)
	if err != nil {
		return err
	}
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(out.Bytes())
	if err != nil {
		return err
	}
	return nil
}
