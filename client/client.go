package client

import (
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
)

type Client struct {
	Jar      *cookiejar.Jar `yaml:"cookies"`
	Username string         `yaml:"username"`
	Password string         `yaml:"password"`
	host     string
	proxy    string
	path     string
	client   *http.Client
}

var Instance *Client

func Init(path, host, proxy string) {
	jar, _ := cookiejar.New(nil)
	c := &Client{Jar: jar, path: path, host: host, proxy: proxy, client: nil}
	if err := c.load(); err != nil {
		color.Red(err.Error())
		color.Green("Create a new session in %v", path)
	}
	Proxy := http.ProxyFromEnvironment
	if len(proxy) > 0 {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			color.Red(err.Error())
			color.Green("use default proxy from environment")
		} else {
			Proxy = http.ProxyURL(proxyURL)
		}
	}
	c.client = &http.Client{Jar: c.Jar, Transport: &http.Transport{Proxy: Proxy}}
	if err := c.save(); err != nil {
		color.Red(err.Error())
	}
	Instance = c
}

func (c *Client) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, c)
}
func (c *Client) save() (err error) {
	data, err := yaml.Marshal(c)
	if err == nil {
		err := os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(c.path, data, 0644)
	}
	if err != nil {
		color.Red("cannot save session to %v\n%v", c.path, err.Error())
	}
	return
}
