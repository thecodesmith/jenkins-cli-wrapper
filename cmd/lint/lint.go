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
	"log"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LintCmd represents the lint command
var LintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a Declarative Jenkinsfile",
	Long: `Lint a Declarative Jenkinsfile.

Automatically lint the Jenkinsfile in the current directory.
Alternatively, provide the path to a Jenkinsfile elsewhere.`,
	Run: func(cmd *cobra.Command, args []string) {
		lint()
	},
}

func init() {
	LintCmd.Flags().StringP("jenkinsfile", "j", "Jenkinsfile", "Path to Jenkinsfile")
	viper.BindPFlag("jenkinsfile", LintCmd.Flags().Lookup("jenkinsfile"))
	// viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
}

func lint() {
	file := viper.Get("jenkinsfile")
	fmt.Printf("Linting %s\n", file)
	fmt.Println("Jenkins host URL:", viper.Get("host"))

	command := fmt.Sprintf("jenkins-cli declarative-linter < %s", file)
	fmt.Println("Running command:", command)

	out, err := exec.Command("sh", "-c", command).CombinedOutput()

	fmt.Printf("Result: %s\n", out)

	if err != nil {
		log.Fatal(err)
	}
}
