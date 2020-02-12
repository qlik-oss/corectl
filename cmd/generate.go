package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

type (
	commandJSON struct {
		Alias       string                 `json:"alias,omitempty"`
		Long        string                 `json:"description,omitempty"`
		Stability   string                 `json:"x-qlik-stability,omitempty"`
		Deprecated  string                 `json:"deprecated,omitempty"`
		Flags       map[string]flagJSON    `json:"flags,omitempty"`
		SubCommands map[string]commandJSON `json:"commands,omitempty"`
	}

	flagJSON struct {
		Shorthand  string `json:"alias,omitempty"`
		Usage      string `json:"description,omitempty"`
		DefValue   string `json:"default,omitempty"`
		Deprecated string `json:"deprecated,omitempty"`
	}

	info struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Version     string `json:"version"`
		License     string `json:"license,omitempty"`
	}

	spec struct {
		Name        string                 `json:"name,omitempty"`
		Info        info                   `json:"info,omitempty"`
		Clispec     string                 `json:"clispec,omitempty"`
		Stability   string                 `json:"x-qlik-stability,omitempty"`
		Flags       map[string]flagJSON    `json:"flags,omitempty"`
		SubCommands map[string]commandJSON `json:"commands,omitempty"`
	}
)

func returnCmdspec(ccmd *cobra.Command) commandJSON {
	ccmdJSON := commandJSON{
		Alias:       returnAlias(ccmd.Aliases),
		Long:        ccmd.Long,
		Deprecated:  ccmd.Deprecated,
		SubCommands: returnCommands(ccmd.Commands()),
		Flags:       returnFlags(ccmd.LocalFlags()),
		Stability:   returnStability(ccmd.Annotations),
	}
	return ccmdJSON
}

func returnAlias(aliases []string) string {
	if len(aliases) != 0 {
		return aliases[0]
	}
	return ""
}

func returnStability(annotations map[string]string) string {
	return annotations["x-qlik-stability"]
}

func returnCommands(commands []*cobra.Command) map[string]commandJSON {
	commadJSON := make(map[string]commandJSON)

	for _, command := range commands {
		commadJSON[strings.Fields(command.Use)[0]] = returnCmdspec(command)
	}
	return commadJSON
}

func returnFlags(flags *pflag.FlagSet) map[string]flagJSON {
	flagsJSON := make(map[string]flagJSON)

	flag := func(f *pflag.Flag) {
		fJSON := flagJSON{
			Shorthand:  f.Shorthand,
			Usage:      f.Usage,
			DefValue:   f.DefValue,
			Deprecated: f.Deprecated,
		}
		flagsJSON[f.Name] = fJSON
	}

	flags.VisitAll(flag)

	return flagsJSON
}

var generateSpecCmd = &cobra.Command{
	Use:    "generate-spec",
	Short:  "Generate API spec based on cobra commands",
	Long:   "Generate API spec docs based on cobra commands",
	Hidden: true,

	Run: func(ccmd *cobra.Command, args []string) {
		fmt.Println("Generating specification")
		var jsonData []byte
		spec := spec{
			Clispec: "0.1.0",
			Name:    ccmd.Root().Use,
			Info: info{
				Title:       "Specification for corectl",
				Description: ccmd.Root().Long,
				Version:     strings.TrimPrefix(version, "v"),
				License:     "MIT",
			},
			SubCommands: returnCommands(ccmd.Root().Commands()),
			Flags:       returnFlags(ccmd.Root().LocalFlags()),
			Stability:   returnStability(ccmd.Root().Annotations),
		}
		jsonData, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		ioutil.WriteFile("./docs/spec.json", jsonData, 0644)
	},
}

var generateDocsCmd = &cobra.Command{
	Use:    "generate-docs",
	Short:  "Generate markdown docs based on cobra commands",
	Long:   "Generate markdown docs based on cobra commands",
	Hidden: true,

	Run: func(ccmd *cobra.Command, args []string) {
		fmt.Println("Generating documentation")
		doc.GenMarkdownTree(ccmd.Root(), "./docs")
	},
}
