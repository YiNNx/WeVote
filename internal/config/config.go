package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var C = &Config{}

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
	Host      string `toml:"host"`
	Port      string `toml:"port"`
}

type Ticket struct {
	Secret string `toml:"secret"`
	Spec   int    `toml:"spec"`
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
	Addr string `toml:"addr"`
}

type Captcha struct {
	RecaptchaSecret string `toml:"recaptcha_secret"`
}

func Init(path string) {
	_, err := toml.DecodeFile(path, C)
	if err != nil {
		panic(err)
	}
	log.Println("Config loaded.")
}
