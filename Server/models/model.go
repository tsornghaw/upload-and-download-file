package models

import (
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Domain     string      `json:"domain,omitempty"`
	Https      HTTPSConfig `json:"https,omitempty"`
	Admin      User        `json:"users,omitempty"`
	Postgresql PQLConfig   `json:"postgresql,omitempty"`
}

type User struct {
	Id       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"password,omitempty"`
	Admin    bool   `jsom:"admin,omitempty" gorm:"column:admin"`
}

type StroeData struct {
	Id          uint      `json:"id" gorm:"primary_key"`
	UploadTime  time.Time `json:"upload_time,omitempty"`
	ShareTime   time.Time `json:"share_time,omitempty"`
	ShareLimit  int       `json:"share_limit,omitempty"`
	FileSize    int64     `json:"file_size,omitempty"`
	FileName    string    `json:"file_name,omitempty"`
	FileType    string    `json:"file_type,omitempty"`
	FileContent string    `json:"file_content,omitempty"`
	DownloadUrl string    `json:"download_url,omitempty"`
}

type DataItem struct {
	Dataname    string `json:"dataname"`
	DownloadURL string `json:"download_url"`
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

type FrontendRequest struct {
	FrontendIDs []uint `json:"frontend_ids"`
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
		Admin: User{
			Id:       0,
			Name:     "admin",
			Email:    "admin@admin.com",
			Password: []byte("adminadmin"),
			Admin:    true,
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
