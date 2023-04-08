package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	LogErr  *log.Logger
	LogWarn *log.Logger
	LogInfo *log.Logger
)

func InitialiseLogging(logFile string) {
	tempDir := os.TempDir()

	_, err := os.Stat(fmt.Sprintf("%s/amnesia", tempDir))
	if os.IsNotExist(err) {
		err = os.Mkdir(fmt.Sprintf("%s/amnesia", tempDir), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/amnesia/%s", tempDir, logFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiOut := io.MultiWriter(os.Stdout, file)

	LogInfo = log.New(multiOut, fmt.Sprintf(
		"(PID: %d) - INFO - ",
		os.Getpid()), log.Ldate|log.Ltime|log.Lshortfile,
	)
	LogWarn = log.New(multiOut, fmt.Sprintf(
		"(PID: %d) - WARNING - ",
		os.Getpid()), log.Ldate|log.Ltime|log.Lshortfile,
	)
	LogErr = log.New(multiOut, fmt.Sprintf(
		"(PID: %d) - ERROR - ",
		os.Getpid()), log.Ldate|log.Ltime|log.Lshortfile,
	)
}
