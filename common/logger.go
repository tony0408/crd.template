package common

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	processid int64
	describe  string
}

func CreateLogger(processid int64, describe string) *Logger {
	return &Logger{processid, describe}
}

func (logger *Logger) Print(message string) {
	location, _ := time.LoadLocation("America/Vancouver")
	t := time.Now().In(location)
	path := fmt.Sprintf("./logs/%02d%02d/", t.Month(), t.Day())
	os.MkdirAll(path, os.ModePerm)

	filename := fmt.Sprintf("%strade-%d-%s.log", path, logger.processid, logger.describe)
	log.Printf("%s", message)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("logger err:%v", err)
	}

	defer f.Close()

	// fmt.Println("Location:", logTime.Location(), ":Time:", logTime)
	logstr := fmt.Sprintf("[%02s/%02d %02d:%02d:%02d] %s\n", t.Month().String(), t.Day(), t.Hour(), t.Minute(), t.Second(), message)
	// fmt.Print(logstr)
	if _, err = f.WriteString(logstr); err != nil {
		log.Printf("logger err:%v", err)
	}
}
