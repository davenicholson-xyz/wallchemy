package wallhaven

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
)

func StartDaemon(app *appcontext.AppContext) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
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
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Daemon shutting down")

	// quit := make(chan os.Signal, 1)
	// <-quit
	// log.Println("Shutting down server...")

}
