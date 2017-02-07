package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func New(aurgoPath string) (Config, error) {
	repoYamlPath := filepath.Join(aurgoPath, "repo.yml")

	repoYamlContents, err := ioutil.ReadFile(repoYamlPath)
	if err != nil {
		return Config{}, err
	}

	var yamlConfig struct {
		Packages []string `yaml:"packages"`
	}
	err = yaml.Unmarshal(repoYamlContents, &yamlConfig)
	if err != nil {
		return Config{}, err
	}

	return Config{
		aurgoPath: aurgoPath,
		packages:  yamlConfig.Packages,
	}, nil
}

type Config struct {
	aurgoPath string
	packages  []string
}

func (c Config) Packages() ([]string, error) {
	return c.packages, nil
}

func (c Config) SourcePath(pkg string) (string, error) {
	sourcePath := filepath.Join(c.aurgoPath, "src", pkg)

	err := os.MkdirAll(sourcePath, 0755)
	if err != nil {
		return "", err
	}

	return sourcePath, nil
}

func (c Config) AurRepoURL(pkg string) string {
	return fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg)
}
