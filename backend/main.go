package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joineroff/social-network/backend/internal/di"
	"go.uber.org/zap"
)

func main() {
	z, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	zap.ReplaceGlobals(z)

	dc := di.New("./config/config.yml")

	ctx, exit := context.WithCancel(context.Background())
	defer exit()

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm

		sctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		dc.Stop(sctx)

		exit()
	}()

	dc.Start()

	<-ctx.Done()
}
