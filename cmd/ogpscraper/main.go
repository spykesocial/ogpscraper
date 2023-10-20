package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yashdiniz/ogpscraper/api"
)

const defaultPort = 8080

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: ogpscraper [-port <port>] [-cache <redis connection sstring>] [-cache-ttl <cache TTL in hours>] [-no-cache]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	port := flag.Int("port", defaultPort, "port to run ogpscraper on (default 8080)")
	cache_disable := flag.Bool("no-cache", false, "disable caching (not recommended in production)")
	cache_uri := flag.String("cache", "localhost:6379", "connection uri for redis cache (defaults to localhost:6379)")
	cache_ttl := flag.Int("cache-ttl", 24, "Number of hours that cache entries are valid for")
	flag.Parse()

	// build configs
	rc := &redis.Options{}
	if cache_uri != nil && *cache_uri != "" {
		rc.Addr = *cache_uri
	}
	log.Printf("Port: %v, cache_disable: %v, cache_ttl: %v, redis_config: %+v", port, cache_disable, cache_ttl, rc)

	// create new server
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	server := &http.Server{
		Addr:    addr,
		Handler: api.NewServer(rc, *cache_disable, *cache_ttl),
	}

	serverCtx := gracefulShutdown(server)
	log.Println("listening on", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panicln(fmt.Errorf("http server terminated: %w", err))
	}

	// Wait for the server to shutdown
	<-serverCtx.Done()
}

// Reference: https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
func gracefulShutdown(server *http.Server) context.Context {
	serverCtx, serverStop := context.WithCancel(context.Background())

	// listen for interrupt signal and gracefully shutdown the server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println(serverCtx, "shutting down server...")
		shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatalln("server shutdown timed out, forcefully terminating...")
			}
			shutdownStop()
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalln(err)
		}
		serverStop()
	}()

	return serverCtx
}
