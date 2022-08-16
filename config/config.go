package config

import (
	"bytes"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type CodeTemplate struct {
	Alias        string   `yaml:"alias"`
	Lang         string   `yaml:"lang"`
	Path         string   `yaml:"path"`
	Suffix       []string `yaml:"suffix"`
	BeforeScript string   `yaml:"before_script"`
	Script       string   `yaml:"script"`
	AfterScript  string   `yaml:"after_script"`
}

type Config struct {
	Template      []CodeTemplate    `yaml:"template"`
	Default       int               `yaml:"default"`
	GenAfterParse bool              `yaml:"gen_after_parse"`
	Host          string            `yaml:"host"`
	Proxy         string            `yaml:"proxy"`
	FolderName    map[string]string `yaml:"folder_name"`
	path          string
}

var Instance *Config

func Init(path string) {
	c := &Config{path: path, Host: "https://atcoder.jp", Proxy: ""}
	if err := c.load(); err != nil {
		color.Red(err.Error())
		color.Green("create a new configuration in %v", path)
	}
	if c.Default < 0 || c.Default >= len(c.Template) {
		c.Default = 0
	}
	if c.FolderName == nil {
		c.FolderName = map[string]string{}
	}
	if _, ok := c.FolderName["root"]; !ok {
		c.FolderName["root"] = "atc"
	}
	err := c.save()
	if err != nil {
		color.Red(err.Error())
		color.Red("create a new configuration in %v, but failed", path)
	}
	Instance = c
}

func (c *Config) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

func (c *Config) save() (err error) {
	var data bytes.Buffer
	encoder := yaml.NewEncoder(&data)
	err = encoder.Encode(c)
	if err == nil {
		err := os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(c.path, data.Bytes(), 0644)
	}
	if err != nil {
		color.Red("cannot save config to %v\n%v", c.path, err.Error())
	}
	return
}
