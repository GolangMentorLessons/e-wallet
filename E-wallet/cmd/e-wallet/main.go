package main

import (
	"E-wallet/cmd/internal"
	"E-wallet/cmd/internal/rest"
	"E-wallet/pkg/repository"
	"os/signal"
	"syscall"

	"os"

	"github.com/sirupsen/logrus"
)

var (
	pgDSN = os.Getenv("PG_DSN")
	port  = os.Getenv("PORT")
)

func main() {
	log := logrus.New()
	pg, err := repository.NewRepo(pgDSN, log)
	if err != nil {
		log.Panicf("Failed to connect database")
	}

	service := internal.NewService(log, pg)

	r := rest.NewRouter(log, service)
	
	go func() {
		if err := r.Run(port); err != nil {
			log.Panicf("Error starting server")
		}
	}()

	//Graceful Shutdown

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	<-sigCh
	pg.Close()
	log.Info("shutting down")

}
