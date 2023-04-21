package configurations

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

type DbConfig struct {
	MeshDbName     string
	BlobDbName     string
	BlobBucketName string
	DbURL          string
}

var IniFilePath string

// ReadDbConfig reads db configuration from config file
func ReadDbConfig(filename string) (*DbConfig, error) {
	return decodeDbConfig(IniFilePath, filename)
}

func decodeDbConfig(path string, filename string) (*DbConfig, error) {
	inifile := filepath.Join(path, filename)
	_, err := os.Stat(inifile)
	log.Println("Configuration: checking existence", inifile)
	if err != nil {
		log.Fatal("Configuration: config file is missing: ", inifile)
	}

	config := &DbConfig{}
	if _, err := toml.DecodeFile(inifile, config); err != nil {
		return nil, fmt.Errorf("could not decode database config: %w", err)
	}
	return config, nil
}
