package yboost

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {
	fil, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Impossibl d'ouvrir le fichier log: %v, err")
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	InfoLogger = log.New(multiWriter, "[INFO] ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(multiWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}