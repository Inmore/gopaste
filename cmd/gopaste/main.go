package main

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	h "github.com/inmore/gopaste/internal/http"
	"github.com/inmore/gopaste/internal/janitor"
	"github.com/inmore/gopaste/internal/storage"
	"github.com/inmore/gopaste/internal/storage/mem"
	"github.com/inmore/gopaste/internal/storage/sqlite"
	"go.uber.org/zap"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "listen address")
		store = flag.String("store", "mem", "mem|sqlite")
		db    = flag.String("dbpath", "data.db", "sqlite file")
	)
	flag.Parse()

	log, _ := zap.NewProduction()
	defer log.Sync()

	var st storage.Storage
	var err error
	switch *store {
	case "mem":
		st = mem.New()
	case "sqlite":
		st, err = sqlite.New(*db)
	default:
		log.Fatal("unknown store", zap.String("store", *store))
	}
	if err != nil {
		log.Fatal("store init", zap.Error(err))
	}
	defer st.Close()

	srv := h.New(log, st)
	server := &http.Server{
		Addr:         *addr,
		Handler:      srv.Routes(),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go janitor.Run(ctx, log, st)

	go func() {
		log.Info("listen", zap.String("addr", *addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http", zap.Error(err))
		}
	}()
	<-ctx.Done()

	shutdownCtx, _ := context.WithTimeout(context.Background(), 6*time.Second)
	_ = server.Shutdown(shutdownCtx)

}
