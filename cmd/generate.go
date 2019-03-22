package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

type (
	commandJSON struct {
		Use        string   `json:"use"`
		Aliases    []string `json:"aliases,omitempty"`
		Short      string   `json:"short,omitempty"`
		Long       string   `json:"long,omitempty"`
		Stability  string   `json:"x-qlik-stability,omitempty"`
		ValidArgs  []string `json:"validArgs,omitempty"`
		Deprecated string   `json:"deprecated,omitempty"`
		// Annotations map[string]string      `json:"annotations,omitempty"`
		Flags       map[string]flagJSON    `json:"flags,omitempty"`
		SubCommands map[string]commandJSON `json:"commands,omitempty"`
	}

	flagJSON struct {
		Name       string `json:"name,omitempty"`
		Shorthand  string `json:"shorthand,omitempty"`
		Usage      string `json:"usage,omitempty"`
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
		Info    info   `json:"info,omitempty"`
		Clispec string `json:"clispec,omitempty"`
		commandJSON
	}
)

func returnCmdspec(ccmd *cobra.Command) commandJSON {
	ccmdJSON := commandJSON{
		Use:        ccmd.Use,
		Aliases:    ccmd.Aliases,
		Short:      ccmd.Short,
		Long:       ccmd.Long,
		ValidArgs:  ccmd.ValidArgs,
		Deprecated: ccmd.Deprecated,
		// Annotations: ccmd.Annotations,
		SubCommands: returnCommands(ccmd.Commands()),
		Flags:       returnFlags(ccmd.LocalFlags()),
		Stability:   returnStability(ccmd.Annotations),
	}
	return ccmdJSON
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
			Name:       f.Name,
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
		var jsonData []byte
		spec := spec{
			Clispec: "0.1.0",
			Info: info{
				Title:       "Specification for corectl",
				Description: "Corectl contains various commands to interact with the Qlik Associative Engine.",
				Version:     version,
				License:     "MIT",
			},
			commandJSON: returnCmdspec(rootCmd),
		}
		jsonData, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(jsonData))
	},
}

var generateDocsCmd = &cobra.Command{
	Use:    "generate-docs",
	Short:  "Generate markdown docs based on cobra commands",
	Long:   "Generate markdown docs based on cobra commands",
	Hidden: true,

	Run: func(ccmd *cobra.Command, args []string) {
		fmt.Println("Generating documentation")
		doc.GenMarkdownTree(rootCmd, "./docs")
	},
}

func init() {
	rootCmd.AddCommand(generateDocsCmd)
	rootCmd.AddCommand(generateSpecCmd)
}
