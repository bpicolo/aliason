package main

import (
	"io/ioutil"
	"os"
	"strings"

	"fmt"

	"strconv"

	"os/exec"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const aliasFilename = ".aliasonrc"
const unaliasEnvVar = "ALIASION_UNENV"
const aliasonInstall = `
function cd() {
    new_directory="$*";
    if [ $# -eq 0 ]; then
        new_directory=${HOME};
    fi;
    builtin cd "${new_directory}" && eval $(aliason env)
}
eval $(aliason env)`

func getRemoveCurrentAliases() string {
	// Unalias will only fail if the alias doesn't currently exist.
	// Should be fine to ignore errors for those cases.
	return fmt.Sprintf("%s %s", os.Getenv(unaliasEnvVar), "&>/dev/null")
}

func sanitizeAlias(a string) string {
	return strconv.Quote(a)
}

func getBuiltins() map[string]bool {
	builtins := map[string]bool{}
	out, err := exec.Command("sh", "-c", "compgen -b").Output()
	if err != nil {
		return builtins
	}

	for _, cmd := range strings.Split(string(out), "\n") {
		builtins[cmd] = true
	}

	return builtins
}

func generateAliasCommand(m map[string]string, builtins map[string]bool) string {
	if len(m) == 0 {
		return ""
	}

	var parts []string
	for k, v := range m {
		if builtins[k] {
			fmt.Fprintf(os.Stderr, "Refusing to create alias for builtin command: %s\n", k)
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, sanitizeAlias(v)))
	}

	return fmt.Sprintf("alias %s", strings.Join(parts, " "))
}

func generateUnaliasCommand(m map[string]string, builtins map[string]bool) string {
	var keys []string
	for k := range m {
		if builtins[k] {
			continue
		}
		keys = append(keys, k)
	}

	return fmt.Sprintf("export %s=\"unalias %s\"", unaliasEnvVar, strings.Join(keys, " "))
}

func unenvAliason() *[]string {
	var commands []string
	if command := getRemoveCurrentAliases(); command != "" {
		commands = append(commands, command)
	}
	commands = append(commands, fmt.Sprintf("unset %s", unaliasEnvVar))

	return &commands
}

func sourceAliasrc() []string {
	var commands []string
	if command := getRemoveCurrentAliases(); command != "" {
		commands = append(commands, command)
		commands = append(commands, fmt.Sprintf("unset %s", unaliasEnvVar))
	}

	if _, err := os.Stat(aliasFilename); os.IsNotExist(err) {
		return commands
	}

	data, err := ioutil.ReadFile(aliasFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .aliasonrc\n")
		return commands
	}

	m := make(map[string]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse .aliasonrc file.\n")
		return commands
	}

	builtins := getBuiltins()
	if command := generateAliasCommand(m, builtins); command != "" {
		commands = append(commands, command)
	}

	if command := generateUnaliasCommand(m, builtins); command != "" {
		commands = append(commands, command)
	}
	return commands
}

func main() {
	app := cli.NewApp()
	app.Name = "aliason"
	app.Usage = "Easily manage project-specific shell aliases"

	app.Commands = []cli.Command{
		{
			Name:  "env",
			Usage: "Source .aliasonrc in the current directory.",
			Action: func(c *cli.Context) error {
				commands := sourceAliasrc()
				if len(commands) > 0 {
					fmt.Println(strings.Join(commands, " &&\n "))
				}
				return nil
			},
		},
		{
			Name:  "unenv",
			Usage: "Unsource the current aliason env.",
			Action: func(c *cli.Context) error {
				commands := unenvAliason()
				if len(*commands) > 0 {
					fmt.Println(strings.Join(*commands, " ;\n "))
				}
				return nil
			},
		},
		{
			Name:  "install",
			Usage: "Generate useful bashrc handlers for enabling aliason",
			Action: func(c *cli.Context) error {
				fmt.Println(aliasonInstall)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
