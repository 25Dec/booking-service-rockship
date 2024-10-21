package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App          `yaml:"app"`
		HTTP         `yaml:"http"`
		Log          `yaml:"logger"`
		PG           `yaml:"postgres"`
		LarkCalendar `yaml:"lark_calendar"`
		LarkBase     `yaml:"lark_base"`
		Edtronaut    `yaml:"edtronaut"`
		Cors         `yaml:"cors"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"false"                 env:"PG_URL"`
	}

	// LarkCalendar -.
	LarkCalendar struct {
		AppID      string `env-required:"true" yaml:"app_id" env:"LARK_APP_ID"`
		AppSecret  string `env-required:"true" yaml:"app_secret" env:"LARK_APP_SECRET"`
		Timezone   string `env-required:"true" yaml:"timezone" env:"LARK_TIMEZONE"`
		CalendarID string `env-required:"true" yaml:"calendar_id" env:"LARK_CALENDAR_ID"`
	}

	// LarkBase -.
	LarkBase struct {
		AppID            string `env-required:"true" yaml:"app_id" env:"LARK_APP_ID"`
		AppSecret        string `env-required:"true" yaml:"app_secret" env:"LARK_APP_SECRET"`
		BaseToken        string `env-required:"true" yaml:"base_token" env:"LARK_BASE_TOKEN"`
		AppointmentTable string `env-required:"true" yaml:"appoitment_table" env:"APPOITMENT_TABLE"`
		CancelTable      string `env-required:"true" yaml:"cancel_table" env:"CANCEL_TABLE"`
	}

	// Edtronaut -.
	Edtronaut struct {
		Domain string `env-required:"true" yaml:"domain" env:"EDTRONAUT_DOMAIN"`
	}

	// Cors -.
	Cors struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	_, er := os.ReadFile("./config/config.yml")
	if er != nil {
		panic(er)
	}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
