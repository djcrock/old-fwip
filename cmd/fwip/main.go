package main

import (
	"context"
	"errors"
	"flag"
	"github.com/djcrock/fwip/internal/repository"
	"github.com/djcrock/fwip/internal/web"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
	"zombiezen.com/go/sqlite/sqlitex"
)

func main() {
	var err error

	bind := flag.String("bind", "", "interface to which the server will bind")
	port := flag.Int("port", 8080, "port on which the server will listen")
	dbStr := flag.String("db", "file:fwip.db", "sqlite3 connection string")
	isVersion := flag.Bool("version", false, "show build and version information")

	flag.Parse()

	info, ok := debug.ReadBuildInfo()
	if ok {
		log.Printf("running fwip build %s %s", info.Main.Version, info.GoVersion)
	}

	if *isVersion {
		if ok {
			log.Printf("build information:\n%v", info)
		} else {
			log.Printf("no version information available")
		}
		return
	}

	var addr = ""
	if *bind != "" {
		ip := net.ParseIP(*bind)
		if ip == nil {
			log.Fatal("invalid ip address provided to --bind")
		}
		addr = ip.String()
	}
	addr += ":" + strconv.Itoa(*port)
	var logAddr = addr
	if *bind == "" {
		logAddr = "localhost" + addr
	}

	dbPool, err := sqlitex.NewPool(*dbStr, sqlitex.PoolOptions{
		Flags:    0,
		PoolSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	repoPool := repository.NewPool(dbPool)
	app := web.NewApp(log.Default(), repoPool)

	srv := &http.Server{
		Addr:    addr,
		Handler: app,
	}

	allConnectionsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 2)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint
		log.Println("received interrupt; shutting down...")

		// Deadline for server shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("error shutting down http server: %v", err)
		}

		log.Println("http server shutdown complete")
		close(allConnectionsClosed)
	}()

	log.Printf("serving fwip at http://%s", logAddr)
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("fatal error in http server: %v", err)
	}

	// Wait for all connections to be closed before terminating
	<-allConnectionsClosed
	log.Println("bye!")
}
