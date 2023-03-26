package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type config struct {
	Database struct {
		DatabaseName string
		Address      string
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
		fmt.Println(err)
		log.Fatalln(err)
	}

	err = viper.Unmarshal(Config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(C)
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
