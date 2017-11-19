package StaticConfiguration

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ReadConfiguration(path string) (config map[string]interface{}, err error) {
	relativePath, _ := filepath.Abs("../" + path)

	var configFile *os.File
	for {
		configFile, err = os.Open(relativePath)

		if err != nil {
			log.Println("Failed to read configuration in path ", relativePath)
			log.Println(err.Error())
			log.Println("Retry in 50 ms")
			time.Sleep(time.Microsecond * 50)
		}else{
			break
		}
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	return

}
