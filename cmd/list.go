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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

func ListContexts() error {
	config, err := ReadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	current, err := config.GetCurrentContext()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, ctx := range config.Contexts {
		if ctx.Name == current.Name {
			fmt.Print("* ")
			color.Green(ctx.Name)
		} else {
			fmt.Printf("  %s\n", ctx.Name)
		}
	}

	return nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured Jenkins contexts",
	Long: `A longer description that spans multiple lines and likely contains examples
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ListContexts()
	},
}

func init() {
	contextCmd.AddCommand(listCmd)
}
