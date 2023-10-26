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
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	// log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	config "github.com/thecodesmith/jenkinsw/pkg/config"
	jenkins "github.com/thecodesmith/jenkinsw/pkg/jenkins"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version info for the Jenkins server, CLI, and wrapper",
	Long:  `Display the version info for the Jenkins server, Jenkins CLI, and Jenkins wrapper CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := printVersion(); err != nil {
			color.Red("Error:", err)
			os.Exit(1)
		}
	},
}

func printVersion() (err error) {
	fmt.Printf("%s version: %s\n", rootCmd.Use, rootCmd.Version)
	fmt.Println("")

	ctx, err := config.GetCurrentContext()

	if err != nil {
		return err
	}

	fmt.Print("Jenkins server: ")
	color.Blue(ctx.Host)
	fmt.Print("  Jenkins server version: ")

	client, err := jenkins.NewClient(&ctx, &ioStreams)
	serverVersion := client.Version()
	fmt.Println(serverVersion)

	// fmt.Print("  Jenkins CLI jar version: ")
	// cli, err := RunJenkinsCli("-version")
	// if err != nil {
	// 	log.Error(err)
	// }

	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
