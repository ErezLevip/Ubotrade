package StaticConfiguration

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func ReadConfiguration(path string) (config map[string]interface{}, err error) {
	relativePath, _ := filepath.Abs("../" + path)
	configFile, err := os.Open(relativePath)
	defer configFile.Close()
	if err != nil {
		log.Panicln("Failed to read configuration in path ", relativePath)
	}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	return

}
