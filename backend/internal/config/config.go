package config

import (
	"github.com/Lavina-Tech-LLC/lavinagopackage/v2/conf"
)

var Confs Conf

type (
	Conf struct {
		DB           gormDB
		Settings     Settings
		JWT          jwtConfig
		AuthProvider AuthProvider
	}

	jwtConfig struct {
		AccessSecret string
	}

	gormDB struct {
		Host     string
		Port     string
		User     string
		Password string
		DbName   string
	}

	Settings struct {
		SrvAddress string
	}

	AuthProvider struct {
		BaseURL      string
		ClientID     string
		ClientSecret string
		RedirectURI  string
	}
)

func init() {
	Confs = conf.Get[Conf]("conf/")
}
