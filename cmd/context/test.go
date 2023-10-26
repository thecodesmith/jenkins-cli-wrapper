/*
Copyright © 2023 Brian Stewart

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

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	config "github.com/thecodesmith/jenkinsw/pkg/config"
	"github.com/thecodesmith/jenkinsw/pkg/jenkins"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test connection to Jenkins",
	Long:  `Test the connection to Jenkins using the URL and credentials from the current context.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := test(); err != nil {
			color.Red("Error: connection failed")
			os.Exit(1)
		}
	},
}

func test() error {
	log.Debug("Testing connection")
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	ctx, err := cfg.GetCurrentContext()
	if err != nil {
		return err
	}

	cli := jenkins.NewJenkinsCli(&ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Connecting to %s as user %s\n", ctx.Host, ctx.Username)
	out, err := cli.RunCommand("who-am-i")
	if err == nil {
		fmt.Println("Success!")
	} else {
		color.Red(string(out))
	}

	return err
}

func init() {
	ContextCmd.AddCommand(testCmd)
}
