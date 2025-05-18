package wallhaven

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
)

func StartDaemon(app *appcontext.AppContext) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
		// var flgValues = make(map[string]any)
		// flgValues["id"] = "4y2wmd"
	})

	go func() {
		pid := os.Getpid()
		err := app.CacheTools.WriteStringToFile("daemon.pid", strconv.Itoa(pid))
		if err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}
		log.Printf("Daemon PID: %d\n", pid)

		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	<-quit
	log.Println("Shutting down server...")

}

func KillDaemon(app *appcontext.AppContext) error {
	data, err := app.CacheTools.ReadLineFromFile("daemon.pid", 1)
	if err != nil {
		log.Fatalf("Failed to read PID file: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return fmt.Errorf("invalid PID in file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	app.CacheTools.DeleteFile("daemon.pid")

	return nil

}
