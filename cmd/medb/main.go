package main

import (
	"encoder-backend/pkg/database"
	"encoder-backend/pkg/encoder"
	"encoder-backend/pkg/http"
	"encoder-backend/pkg/manager"
	"encoder-backend/pkg/models"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	Version string
	Build   string
)

func main() {

	measure := time.Now()

	log.WithFields(log.Fields{
		"version": Version,
		"build":   Build,
	}).Info("runtime: starting")

	preload()

	web := http.New()
	manage := manager.New()
	encode := encoder.New()

	wait := &sync.WaitGroup{}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	wait.Add(1)
	go func() {
		defer wait.Done()
		<-sig
		log.Info("runtime: shutdown requested")
	}()

	log.WithField("duration", time.Since(measure)).Info("runtime: started")

	wait.Wait()

	web.Close()
	manage.Close()
	encode.Close()
}

func preload() {

	database.Migrate(
		models.Encode{},
		models.File{},
		models.Path{},
		models.QualityProfile{},
		models.Setting{},
		models.Revision{},
	)
}
