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
	"path/filepath"

	"github.com/ghodss/yaml"
)

type Config struct {
	Contexts       []Context `json:"contexts"`
	CurrentContext string    `json:"current-context"`
}

type Context struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	ApiToken string `json:"apiToken"`
}

const ConfigDir = ".jenkinsw"
const ConfigFile = "config"

func ReadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	dir := filepath.Join(homeDir, ConfigDir)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return Config{}, err
	}

	f := filepath.Join(homeDir, ConfigDir, ConfigFile)

	if _, err := os.Stat(f); os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist, creating it\n", f)
		f2, err := os.Create(f)
		if err != nil {
			return Config{}, err
		}

		defer f2.Close()

		err = os.Chmod(f, 0600)
		if err != nil {
			return Config{}, err
		}
	}

	y, err := os.ReadFile(f)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(y, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) Save() error {
	y, err := yaml.Marshal(c)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	dir := filepath.Join(homeDir, ConfigDir)

	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return fmt.Errorf("Error: %s\n", err)
	}

	configFile := filepath.Join(dir, ConfigFile)

	return os.WriteFile(configFile, y, 0600)
}

func GetCurrentContext() (Context, error) {
	c, err := ReadConfig()
	if err != nil {
		return Context{}, err
	}

	return c.GetCurrentContext()
}

func (c Config) GetCurrentContext() (Context, error) {
	for _, ctx := range c.Contexts {
		if ctx.Name == c.CurrentContext {
			return ctx, nil
		}
	}

	return Context{}, fmt.Errorf("Context named '%s' not found", c.CurrentContext)
}

func (c Config) GetContext(name string) (Context, error) {
	for _, ctx := range c.Contexts {
		if ctx.Name == name {
			return ctx, nil
		}
	}

	return Context{}, fmt.Errorf("Context named '%s' not found", c.CurrentContext)
}

func (c Config) AddContext(context Context) error {
	for _, ctx := range c.Contexts {
		if ctx.Name == context.Name {
			return fmt.Errorf("Context named '%s' already exists", context.Name)
		}
	}

	c.Contexts = append(c.Contexts, context)

	if len(c.CurrentContext) < 1 {
		c.CurrentContext = context.Name
	}

	return nil
}

func (c Config) UseContext(name string) error {
	if !c.IsExistingContext(name) {
		return fmt.Errorf("No context named '%s'. Use 'jenkinsw context list' to view available contexts.", name)
	}

	c.CurrentContext = name

	return c.Save()
}

func (c Config) IsExistingContext(name string) bool {
	exists := false

	for _, ctx := range c.Contexts {
		if ctx.Name == name {
			exists = true
		}
	}

	return exists
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ConfigDir), nil
}

func (c Config) GetConfigFile() (string, error) {
	path, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, ConfigFile), nil
}

func (c Context) GetAuthFile() (string, error) {
	dir, err := c.GetContextDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, ".auth"), nil
}

func (c Context) GetContextDir() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "context", c.Name), nil
}

func (c Context) SaveAuthFile() error {
	dir, err := c.GetContextDir()
	if err != nil {
		return err
	}

	file, err := c.GetAuthFile()
	if err != nil {
		return err
	}

	// Create context directory in case it does not exist
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	auth := fmt.Sprintf("%s:%s", c.Username, c.ApiToken)

	return os.WriteFile(file, []byte(auth), 0600)
}
