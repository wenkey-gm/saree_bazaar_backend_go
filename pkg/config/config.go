package config

import (
	"path"
	"path/filepath"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type config struct {
	Database struct {
		DatabaseName   string
		Address        string
	}
	Server struct {
		Port string
	}
}

var C config

func ReadConfig() {
	Config := &C

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(rootDir(), "config"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(Config)

	if err != nil {
		panic(err)
	}

	spew.Dump(Config)
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(d)
}
