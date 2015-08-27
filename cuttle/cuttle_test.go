package main

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestMain(m *testing.M) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)

	os.Exit(m.Run())
}
