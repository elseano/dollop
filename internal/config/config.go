package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MessageField   string
	Groups         []GroupSpec
	TimestampField string
}

type GroupSpec struct {
	ValueField string
	TitleField string
	Name       string
}

func Get() (config Config) {
	config = Config{
		MessageField:   "msg",
		TimestampField: "time",
		Groups: []GroupSpec{
			{ValueField: "request_id", TitleField: "msg", Name: "Request"},
			{ValueField: "category", TitleField: "category", Name: "Category"},
		},
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
