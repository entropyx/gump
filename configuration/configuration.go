package configuration

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Gump    string              `yaml:"gump_version,omitempty"`
	Version string              `yaml:"version,omitempty"`
	Files   []FileConfiguration `yaml:"files,omitempty"`
}

type FileConfiguration struct {
	Path   string   `yaml:"path,omitempty"`
	Keys   []string `yaml:"keys,omitempty"`
	Prefix string   `yaml:"prefix,omitempty"`
}

func (c *Configuration) Write(path string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 644)
	if err != nil {
		return err
	}
	return nil
}
