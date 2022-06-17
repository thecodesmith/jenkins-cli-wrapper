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
package lint

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thecodesmith/jenkinsw/cmd/context"
)

var debugMode bool

// LintCmd represents the lint command
var LintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a Declarative Jenkinsfile",
	Long: `Lint a Declarative Jenkinsfile.

Automatically lint the Jenkinsfile in the current directory.
Alternatively, provide the path to a Jenkinsfile elsewhere.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := lint(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	LintCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug output")
	LintCmd.Flags().StringP("jenkinsfile", "j", "Jenkinsfile", "Path to Jenkinsfile")
	viper.BindPFlag("jenkinsfile", LintCmd.Flags().Lookup("jenkinsfile"))
	// viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
}

func debug(a ...any) {
	if debugMode {
		fmt.Println(a)
	}
}

func lint() error {
	jenkinsfile := viper.Get("jenkinsfile")
	debug("Linting", jenkinsfile)

	config, err := context.ReadConfig()
	if err != nil {
		return err
	}

	ctx, err := config.GetCurrentContext()
	if err != nil {
		return err
	}

	cli, err := ctx.GetCliPath()
	if err != nil {
		return err
	}

	authFile, err := ctx.GetAuthFile()
	if err != nil {
		return err
	}

	if _, err := os.Stat(authFile); err != nil {
		return fmt.Errorf("Authentication file not present for context '%s'. Please run 'jenkinsw context add' again.", ctx.Name)
	}

	command := fmt.Sprintf("java -jar '%s' -s '%s' -auth '@%s' -webSocket declarative-linter < '%s'", cli, ctx.Host, authFile, jenkinsfile)

	cmd := exec.Command("sh", "-c", command)

	debug("Running command:", command)
	out, err := cmd.CombinedOutput()

	debug("Result:", string(out))

	return err
}
