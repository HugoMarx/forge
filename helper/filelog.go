package helper

import (
	"log"
	"os"
)

const logDir = "./log/debug.log"

func LogToDebug(output string) {
	logFile, err := os.OpenFile(logDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Panic(err.Error())
	}

	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)
	log.Println(output)
}
