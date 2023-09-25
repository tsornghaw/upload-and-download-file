package models

import (
	"os"
	"path/filepath"
)

type Config struct {
	Domain     string      `json:"domain,omitempty"`
	Https      HTTPSConfig `json:"https,omitempty"`
	Postgresql PQLConfig   `json:"postgresql,omitempty"`
}

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}

type HTTPSConfig struct {
	ListenIp   string    `json:"listenIP,omitempty"`
	ListenPort uint16    `json:"listenPort,omitempty"`
	TLS        TLSConfig `json:"tls,omitempty"`
}

type TLSConfig struct {
	Cert string `json:"cert,omitempty"`
	Key  string `json:"key,omitempty"`
}

type PQLConfig struct {
	UserName     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	DatabaseName string `json:"database_name,omitempty"`
	DatabaseType string `json:"database_type,omitempty"`
}

var (
	currentDir, _ = os.Getwd()
	parentDir     = filepath.Dir(currentDir)
	DefaultConfig = Config{
		Domain: "localhost",
		Https: HTTPSConfig{
			ListenIp:   "0.0.0.0",
			ListenPort: 4443,
			TLS: TLSConfig{
				Cert: filepath.Join(parentDir, "certs", "public_key.pem"),
				Key:  filepath.Join(parentDir, "certs", "private_key.pem"),
			},
		},
		Postgresql: PQLConfig{
			UserName:     "postgres",
			Password:     "mysecretpassword",
			Host:         "localhost",
			Port:         5432,
			DatabaseName: "postgres",
			DatabaseType: "postgres",
		},
	}
)
