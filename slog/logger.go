package slog

import (
	"fmt"
	"log"
	"net/http"
)

type ID struct{}

const (
	StyleBlue      = "\033[94m"
	StyleCyan      = "\033[96m"
	StyleGreen     = "\033[92m"
	StyleYellow    = "\033[93m"
	StyleRed       = "\033[91m"
	StyleEnd       = "\033[0m"
	StyleBold      = "\033[1m"
	StyleUnderline = "\033[4m"
)

func Style(message, color string) string {
	return fmt.Sprintf("%s%s%s", color, message, StyleEnd)
}

func corrID(r *http.Request) string {
	corrId := r.Context().Value(ID{})
	return fmt.Sprintf("%s[%s]%s", StyleCyan, corrId, StyleEnd)
}

func Info(logger *log.Logger, r *http.Request, message string, v ...any) {
	logMessage := fmt.Sprintf(message, v...)
	logger.Printf("%s %s", corrID(r), logMessage)
}

func Warning(logger *log.Logger, r *http.Request, message string, v ...any) {
	logMessage := fmt.Sprintf(message, v...)
	logger.Printf("%s %s %s", corrID(r), Style("WARNING:", StyleYellow), logMessage)
}

func Error(logger *log.Logger, r *http.Request, message string, v ...any) {
	logMessage := fmt.Sprintf(message, v...)
	logger.Printf("%s %s %s", corrID(r), Style("ERROR:", StyleRed), logMessage)
}
