package Logger

import (
	"os"
	"path/filepath"
	"log"
	"io"
)

func SetGlobalLogger()  {
	wd, err := os.Getwd()
	if (err != nil) {
		log.Fatal(err.Error())
	}
	_, fileLocation := filepath.Split(wd)
	filePath :=  fileLocation + ".txt"

	logFile ,err := os.OpenFile(filePath,os.O_CREATE | os.O_APPEND | os.O_RDWR,0666)
	if(err != nil){
		log.Panic(err.Error())
	}
	mw := io.MultiWriter(os.Stdout,logFile)
	log.SetOutput(mw)
}
