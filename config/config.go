package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Config struct {
	Http struct {
		Host            string        `json:"host" validate:"required"`
		Port            uint16        `json:"port"`
		TlsCertFile     string        `json:"tls_cert_file" validate:"required,file"`
		TlsKeyFile      string        `json:"tls_key_file" validate:"required,file"`
		PublicDir       string        `json:"public_dir" validate:"required,dir"`
		ShutdownTimeout time.Duration `json:"-"`
	} `json:"http"`

	Storage struct {
		Timeout int `json:"timeout" validate:"min=-1"`

		Attachments struct {
			FilePattern string `json:"file_pattern" validate:"contains=*"`
			MaxSizeMb   int    `json:"max_size_mb" validate:"min=-1"`
		} `json:"attachments"`

		Bunt *struct {
			File string `json:"file" validate:"required"`
		} `json:"bunt" validate:"required_without_all=Mongo"`

		Mongo *struct {
			Host     string `json:"host" validate:"required"`
			Port     uint16 `json:"port"`
			Rs       string `json:"rs" validate:"required"`
			Db       string `json:"db" validate:"required"`
			User     string `json:"user" validate:"required"`
			PassFile string `json:"pass_file" validate:"required,file"`
		} `json:"mongo" validate:"required_without_all=Bunt"`
	} `json:"storage" validate:"dive"`
}

var validate = validator.New()

func New(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config: %v", err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("decode config: %v", err)
	}
	if err := validate.Struct(&c); err != nil {
		return nil, err
	}

	// hidden constants
	c.Http.ShutdownTimeout = 10 * time.Second

	return &c, err
}

func ReadFile(filename, defVal string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("read file from config: %v\n", err)
		return defVal
	}
	return strings.TrimRight(string(b), "\r\n")
}
