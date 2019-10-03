package configurations

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type DbConfig struct {
	Dbname string
	Dburl string
}

// Reads info from config file
func ReadConfig(filename string) DbConfig {
	var configfile = filename
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config DbConfig
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}