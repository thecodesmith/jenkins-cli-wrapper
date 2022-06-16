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
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var contextAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Jenkins context",
	Long:  `Add a new Jenkins context configuration, by providing flag parameters or interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := prompt("Context name: ")
		host, _ := prompt("Jenkins URL: ")

		fmt.Println("name:", name)

		jenkinsJarUrl := fmt.Sprintf("%s/jnlpJars/jenkins-cli.jar", host)
		fmt.Printf("Downloading Jenkins CLI from %s\n", jenkinsJarUrl)

		// Download CLI jar file from Jenkins host
		err := download("cli.jar", jenkinsJarUrl)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		// Prompt for username + API token, store in secure file
		username, _ := prompt("Jenkins username: ")
		apiToken, _ := prompt("Jenkins API token: ")

		config, err := ReadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		context := Context{name, host, username, apiToken}
		config.Contexts = append(config.Contexts, context)
		config.CurrentContext = name

		y, err := yaml.Marshal(config)

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		dir := filepath.Join(homeDir, ConfigDir)

		err = os.MkdirAll(dir, 0700)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		configFile := filepath.Join(dir, ConfigFile)
		err = os.WriteFile(configFile, y, 0600)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	contextCmd.AddCommand(contextAddCmd)
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

func download(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
