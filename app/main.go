package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newRouter() *httprouter.Router { //will return the multiplexer
	mux := httprouter.New()
	mux.GET("/youtoube/channel/stats", getChannelStats())
	return mux
}

func getChannelStats() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		write, err := writer.Write([]byte("response"))
		if err != nil {
			return
		}
		fmt.Println(write)
	}
}

func main() {
	server := &http.Server{
		Addr:    ":8090",
		Handler: newRouter(),
	}
	idleConsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("Service interrupt received")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFunc()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown err:%w", err)

		}
		log.Println("shutwdoen complete")
		close(idleConsClosed)
	}()
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Http server failed to start: %v", err)
		}
	}
	<-idleConsClosed
	log.Println("service stoped")
}
