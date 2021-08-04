package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"handler/config"
	"handler/handlers"
	"handler/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	writeTimeout    = time.Second * 15
	readTimeout     = time.Second * 15
	idleTimeout     = time.Second * 60
	gracefulTimeout = time.Second * 15
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Panicln(err.Error())
	}

	// Start metrics server in a goroutine
	go runMetricsServer(cfg.MetricsPort)

	routes := mux.NewRouter()

	// define you handlers here
	routes.HandleFunc("/", handlers.MakeInfoHandler()).Methods("POST")

	// do not change these handlers
	routes.HandleFunc("/_/ready", handlers.MakeHealthHandler()).Methods("GET")
	routes.HandleFunc("/_/health", handlers.MakeHealthHandler()).Methods("GET")

	// Setup middleware handler with metrics
	middlewareOptions := metrics.MakeMiddlewareOptions()
	routes.Use(middlewareOptions.MiddlewareHandler)

	// На все остальные end-point выдаем статус 501
	routes.PathPrefix("/").Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Not Implemented request - %v (%v)", r.RequestURI, r.Method)
			w.WriteHeader(http.StatusNotImplemented)
		}))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ListenPort),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      routes, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panicf("can't start listen: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	_ = srv.Shutdown(ctx)

	log.Panicln("shutting down")
}

// Listen on a separate HTTP port for Prometheus metrics to keep this accessible from
// the internal network only.
func runMetricsServer(port int) {
	router := mux.NewRouter()
	router.Handle("/metrics", metrics.PrometheusHandler())

	log.Printf("start metrics server on \"%d\"", port)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        router,
	}

	log.Panicln(s.ListenAndServe().Error())
}