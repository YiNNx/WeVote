package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

var C *Config

type Config struct {
	Server   Server   `toml:"server"`
	Postgres Postgres `toml:"postgres"`
	Redis    Redis    `toml:"redis"`
	Log      Log      `toml:"log"`
	Vote     Vote     `toml:"vote"`
	Ticket   Ticket   `toml:"ticket"`
	Captcha  Captcha  `toml:"captcha"`
}

type Server struct {
	DebugMode bool   `toml:"debug_mode"`
	Addr      string `toml:"addr"`
}

type Vote struct {
	CacheTTL        time.Duration `toml:"cache_ttl"`
	WriteBackPeriod duration      `toml:"write_back_period"`
}

type Ticket struct {
	Secret     string   `toml:"secret"`
	TTL        duration `toml:"ttl"`
	UsageLimit int      `toml:"usage_limit"`
}

type Log struct {
	Path string `toml:"path"`
}

type Postgres struct {
	DSN string `toml:"dsn"`
}

type Redis struct {
	Addrs []string `toml:"addr"`
}

type Captcha struct {
	Open            bool   `toml:"open"`
	RecaptchaSecret string `toml:"recaptcha_secret"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func Init(path string) error {
	var config *Config
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return err
	}
	log.Println("Config loaded.")
	C = config
	return nil
}
