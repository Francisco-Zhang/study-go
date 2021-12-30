package main

import (
	"go.uber.org/zap"
)

func main() {
	// bb.Nn()
	log, _ := zap.NewProduction()
	log.Warn("warn test")
}
