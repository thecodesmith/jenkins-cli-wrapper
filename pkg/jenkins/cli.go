package jenkins

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	config "github.com/thecodesmith/jenkinsw/pkg/config"
)

type JenkinsCli struct {
	ctx       *config.Context
}

func NewJenkinsCli(ctx *config.Context) JenkinsCli {
	return JenkinsCli{ctx: ctx}
}

func (c JenkinsCli) GetCliDir() (string, error) {
	r := strings.NewReplacer("://", "_", "/", "_", ":", "_")
	hostDir := r.Replace(c.ctx.Host)

	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "cli", hostDir), nil
}

func (c JenkinsCli) GetCliPath() (string, error) {
	dir, err := c.GetCliDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "cli.jar"), nil
}

func (c JenkinsCli) DownloadCliJar() error {
	jenkinsJarUrl := fmt.Sprintf("%s/jnlpJars/jenkins-cli.jar", c.ctx.Host)

	dir, err := c.GetCliDir()
	if err != nil {
			return err
	}

	path, _ := c.GetCliPath()

	fmt.Printf("Downloading Jenkins CLI from %s to %s\n", jenkinsJarUrl, path)

	// Create CLI directory
	if err = os.MkdirAll(dir, 0700); err != nil {
			return err
	}

	// Download CLI jar file from Jenkins host
	return Download(path, jenkinsJarUrl)
}

func Download(filepath string, url string) (err error) {
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

func (c JenkinsCli) RunCommand(subcommand string) (out []byte, err error) {
	cli, err := c.GetCliPath()
	if err != nil {
		return nil, err
	}

	authFile, err := c.ctx.GetAuthFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(authFile); err != nil {
		return nil, fmt.Errorf("Authentication file not present for context '%s'. Please run 'jenkinsw context add' again.", c.ctx.Name)
	}

	command := fmt.Sprintf("java -jar '%s' -s '%s' -auth '@%s' -webSocket %s", cli, c.ctx.Host, authFile, subcommand)

	cmd := exec.Command("sh", "-c", command)

	log.Debug("Running command:", command)
	out, err = cmd.CombinedOutput()

	return out, err
}

func (c JenkinsCli) Version() (version string, err error) {
	path, err := c.GetCliPath()
	if err != nil {
		return "", err
	}

	archive, err := zip.OpenReader(path)
    if err != nil {
        return "", err
    }
    defer archive.Close()

	for _, f := range archive.File {
		if f.Name != "META-INF/MANIFEST.MF" {
			continue
		}

        manifestFile, err := f.Open()
        if err != nil {
            return "", err
        }
        defer manifestFile.Close()

		content := new(strings.Builder)
        if _, err := io.Copy(content, manifestFile); err != nil {
            return "", err
        }

		version := "unknown"
		scanner := bufio.NewScanner(strings.NewReader(content.String()))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Jenkins-CLI-Version:") {
				words := strings.Fields(line)
				if len(words) == 2 {
					version = words[1]
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}

        return version, nil
    }

	return "unknown", nil
}
