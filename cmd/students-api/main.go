package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"students/internal/config"
	"students/internal/http/handlers/student"
	"syscall"
	"time"
)

func main(){
	// load config
	cfg:= config.MustLoad()
	// database setup
	// router setup
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())
	// server setup

	server := http.Server {
		Addr: cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("addres", cfg.HTTPServer.Addr))

	done:= make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func(){
		err:= server.ListenAndServe()
		if err != nil{
			log.Fatalf("failed to start server, %s", err.Error())
		}
	}()

	<- done

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err:= server.Shutdown(ctx); err != nil{
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}	
}