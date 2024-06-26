package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
)

var Env EnvType

type EnvType struct {
	Port string `default:"8080"`
	App  struct {
		Secret string `required:"true"`
		URL    string `required:"true"`
	}
	DB struct {
		Socket   string
		Host     string
		Port     uint
		User     string `required:"true"`
		Password string `required:"true"`
		Name     string `required:"true"`
	}
}

func fileExists(path string) bool {
	if f, err := os.Stat(path); os.IsNotExist(err) || f.IsDir() {
		return false
	} else {
		return true
	}
}

func init() {
	dotenvPath := "./.env"
	if dotenvPathEnv := os.Getenv("DOTENV_PATH"); dotenvPathEnv != "" {
		dotenvPath = dotenvPathEnv
	}
	if fileExists(dotenvPath) {
		err := godotenv.Load(dotenvPath)
		if err != nil {
			panic(err)
		}
	}

	err := envconfig.Process("", &Env)
	if err != nil {
		panic(err)
	}
}
