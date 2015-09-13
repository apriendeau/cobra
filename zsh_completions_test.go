package cobra

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenZshCompletion(t *testing.T) {
	root := &Command{
		Use:   "test [subcommand]",
		Short: "Testing Zsh Generator",
		Run:   func(cmd *Command, args []string) {},
	}
	subCommand := &Command{
		Use:   "subcommand [args]",
		Short: "Adding an short commands",
	}
	root.AddCommand(subCommand)
	out := new(bytes.Buffer)
	err := root.GenZshCompletion(out)
	if err != nil {
		t.Error(err)
	}
}

func TestGenZshCompletionFile(t *testing.T) {
	dir, err := ioutil.TempDir(".", "cobra-tmp")
	if err != nil {
		t.Error(err)
	}

	root := &Command{
		Use:   "test [subcommand]",
		Short: "Testing Zsh Generator",
		Run:   func(cmd *Command, args []string) {},
	}
	subCommand := &Command{
		Use:   "subcommand [args]",
		Short: "Adding an short commands",
	}

	root.AddCommand(subCommand)
	err = root.GenZshCompletionFile(dir + "/_sample")
	if err != nil {
		t.Error(err)
	}
	o, err := ioutil.ReadFile(dir + "/_sample")
	if err != nil {
		t.Error(err)
	}
	if string(o) == "" {
		t.Error(err)
	}
	err = os.RemoveAll(dir)
	if err != nil {
		t.Error(err)
	}
}
