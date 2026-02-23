package config

import (
	"github.com/Lavina-Tech-LLC/lavinagopackage/v2/conf"
)

var Confs Conf

type (
	Conf struct {
		DB       gormDB
		Settings Settings
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
		JWTSecret  string
	}
)

func Init() {
	Confs = conf.Get[Conf]("conf/")
}
