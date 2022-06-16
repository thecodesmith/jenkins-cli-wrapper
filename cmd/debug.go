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
	"strings"
	"path/filepath"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func PrintConfigDetails() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f := filepath.Join(homeDir, ConfigDir, ConfigFile)
	fmt.Println("Config files:")
	fmt.Println("-", f)

	y, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	fmt.Println()

	color.Set(color.FgWhite, color.Bold)
	fmt.Printf("Contents of %s:\n", f)
	color.Unset()
	fmt.Println()
	PrintYaml(string(y))

	return nil
}

func PrintYaml(s string) error {
	// Indent string block
	s = fmt.Sprintf("    %s\n", strings.Replace(s, "\n", "\n    ", -1))

	return quick.Highlight(os.Stdout, s, "yaml", "terminal256", "github")
}

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := PrintConfigDetails()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	contextCmd.AddCommand(debugCmd)
}
