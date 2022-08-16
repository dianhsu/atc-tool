package client

import (
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
)

type Client struct {
	Jar      *cookiejar.Jar `json:"cookies"`
	Username string         `json:"username"`
	Password string         `json:"password"`
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
}

func (c *Client) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c)
}
func (c *Client) save() (err error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		err = os.WriteFile(c.path, data, 0644)
	}
	if err != nil {
		color.Red("Cannot save session to %v\n%v", c.path, err.Error())
	}
	return
}
