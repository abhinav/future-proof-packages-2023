// serve watches for changes to adoc files with entr,
// and rebuilds the site when they change.
//
// It serves the contents of the public directory on a fixed port.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
)

var _port = flag.Int("port", 8080, "port to serve on")

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go serveHTTP(ctx, &wg, *_port)
	go watchFiles(ctx, &wg)

	<-ctx.Done()
}

func serveHTTP(ctx context.Context, wg *sync.WaitGroup, port int) {
	defer wg.Done()

	http.Handle("/", http.FileServer(http.Dir(".")))
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	context.AfterFunc(ctx, func() {
		log.Printf("shutting down")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown: %v", err)
		}
	})

	log.Printf("serving on http://localhost:%d", port)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func watchFiles(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if err := ctx.Err(); err != nil {
			return
		}

		if err := watchFilesOnce(ctx); err != nil {
			// errors from entr are expected
			log.Printf("watch files: %v", err)
		}
	}
}

// uses entr to watch for changes to adoc files
// and rebuilds the site when they change.
//
// If a new adoc file is added, entr will fail (-d flag).
func watchFilesOnce(ctx context.Context) error {
	adocs, err := filepath.Glob("*.adoc")
	if err != nil {
		return fmt.Errorf("glob: %w", err)
	}

	cmd := exec.CommandContext(ctx, "entr", "-d", "make", "build-dev")
	cmd.Stdin = strings.NewReader(strings.Join(adocs, "\n"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Cancel = func() error {
		log.Printf("stopping watch")
		return cmd.Process.Signal(os.Interrupt)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("entr: %w", err)
	}

	return nil
}
