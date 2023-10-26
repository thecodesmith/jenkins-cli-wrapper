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

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [context]",
	Short: "Show details of a specific Jenkins context (default: current context)",
	Long:  `Show details of the specified Jenkins context. Defaults to showing the current context if none is specified.`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		var context config.Context

		if len(args) == 0 {
			context, err = cfg.GetCurrentContext()
		} else if len(args) == 1 {
			context, err = cfg.GetContext(args[0])
		}

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Context:", context.Name)
		fmt.Println("Jenkins URL:", context.Host)
		fmt.Println("Username:", context.Username)
		fmt.Println("API Token:", context.ApiToken)
	},
}

func init() {
	ContextCmd.AddCommand(showCmd)
}
