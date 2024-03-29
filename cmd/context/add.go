/*
Copyright © 2022 Brian Stewart

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	config "github.com/thecodesmith/jenkinsw/pkg/config"
)

var contextAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Jenkins context",
	Long:  `Add a new Jenkins context configuration, by providing flag parameters or interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := prompt("Context name: ")
		host, _ := prompt("Jenkins URL: ")

		// Prompt for username + API token, store in secure file
		username, _ := prompt("Jenkins username: ")
		apiToken, _ := prompt("Jenkins API token: ")

		cfg, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		context := config.Context{name, host, username, apiToken}
		cfg.Contexts = append(cfg.Contexts, context)
		cfg.CurrentContext = name

		err = cfg.Save()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		err = context.SaveAuthFile()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	ContextCmd.AddCommand(contextAddCmd)
}

func prompt(text string) (string, error) {
	fmt.Print(text)
	var result string
	_, err := fmt.Scanln(&result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return result, nil
}
