package main

import (
	"io/ioutil"
	"os"
	"strings"

	"fmt"

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
	return os.Getenv(unaliasEnvVar)
}

func generateAliasCommand(m *map[string]string) string {
	if len(*m) == 0 {
		return ""
	}

	var parts []string
	for k, v := range *m {
		parts = append(parts, fmt.Sprintf("%s=\"%s\"", k, v))
	}

	return fmt.Sprintf("alias %s", strings.Join(parts, " "))
}

func generateUnaliasCommand(m *map[string]string) string {
	var keys []string
	for k := range *m {
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

func sourceAliasrc() *[]string {
	var commands []string
	if command := getRemoveCurrentAliases(); command != "" {
		commands = append(commands, command)
	}

	if _, err := os.Stat(aliasFilename); os.IsNotExist(err) {
		commands = append(commands, fmt.Sprintf("unset %s", unaliasEnvVar))
		return &commands
	}

	data, err := ioutil.ReadFile(aliasFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .aliasonrc")
	}

	m := make(map[string]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse .aliasonrc file.")
	}

	if command := generateAliasCommand(&m); command != "" {
		commands = append(commands, command)
	}

	if command := generateUnaliasCommand(&m); command != "" {
		commands = append(commands, command)
	}

	return &commands
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
				if len(*commands) > 0 {
					fmt.Println(strings.Join(*commands, " &&\n "))
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
