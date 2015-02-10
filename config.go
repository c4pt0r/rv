package rv

import (
	"encoding/json"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type VHost struct {
	hostNamePattern string   `toml:"hostname"`
	Upstreams       []string `toml:"upstreams"`
	Static          string   `toml:"static"`
}

type Config struct {
	VHost []VHost `toml:"vhost"`
	Addr  string  `toml:"addr"`
}

func (c Config) String() string {
	// output as json, for debug
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

func loadConfig(filename string) (Config, error) {
	conf := Config{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(data), &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
