package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   Server   `toml:"server"`
	Postgres Postgres `toml:"postgres"`
	Redis    Redis    `toml:"redis"`
	Log      Log      `toml:"log"`
	Ticket   Ticket   `toml:"ticket"`
	Captcha  Captcha  `toml:"captcha"`
}

type Server struct {
	DebugMode bool   `toml:"debug_mode"`
	Addr      string `toml:"addr"`
}

type Ticket struct {
	Secret     string `toml:"secret"`
	Spec       int    `toml:"spec"`
	UpperLimit int    `toml:"upper_limit"`
}

type Log struct {
	Path string `toml:"path"`
}

type Postgres struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
}

type Redis struct {
	Addrs []string `toml:"addr"`
}

type Captcha struct {
	RecaptchaSecret string `toml:"recaptcha_secret"`
}

func Init(path string) *Config {
	var config *Config
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		panic(err)
	}
	log.Println("Config loaded.")
	return config
}
