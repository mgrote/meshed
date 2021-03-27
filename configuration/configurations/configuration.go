package configurations

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

type DbConfig struct {
	Dbname string
	Dburl string
	Bucketname string
}

var IniFilePath string

// Reads db configuration from config file
func ReadDbConfig(filename string) DbConfig {
	return decodeDbConfig(IniFilePath, filename)
}

func decodeDbConfig(path string, filename string) DbConfig {
	inifile := filepath.Join(path, filename)
	_, err := os.Stat(inifile)
	if err != nil {
		log.Fatal("Config file is missing: ", inifile)
	}

	var config DbConfig
	if _, err := toml.DecodeFile(inifile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
