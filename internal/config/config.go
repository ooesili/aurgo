package config

import (
	"fmt"
	"io/ioutil"
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

func (c Config) Packages() []string {
	return c.packages
}

func (c Config) SourcePath(pkg string) string {
	return filepath.Join(c.aurgoPath, "src", pkg)
}

func (c Config) AurRepoURL(pkg string) string {
	return fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg)
}

func (c Config) SourceBase() string {
	return filepath.Join(c.aurgoPath, "src")
}
