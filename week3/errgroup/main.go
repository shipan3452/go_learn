package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancle := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signalChannel
		cancle()
		return fmt.Errorf("received exit signal: %s", sig)
	})

	g.Go(func() error {
		<-gctx.Done()
		if err := srv.Shutdown(context.TODO()); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print("listen: %+V\n", err)
			return err
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		log.Fatalf("error group wait error :%+v", err)
	} else {
		fmt.Println("finished clean")
	}
}
