package wallhaven

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func StartDaemon(app *appcontext.AppContext) {

	go func() {
		pid := os.Getpid()
		err := app.CacheTools.WriteStringToFile("daemon.pid", strconv.Itoa(pid))
		if err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}
		log.Printf("Daemon PID: %d\n", pid)

		defer func() {
			app.CacheTools.DeleteFile("daemon.pid")
		}()

		mux := http.NewServeMux()

		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			fmt.Fprintf(w, "{\"ping\":\"pong\"}")
		})

		mux.HandleFunc("GET /wp/{id}", func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)

			id := r.PathValue("id")

			cmd := exec.Command("wallchemy", "-id", id)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()

			if err := cmd.Start(); err != nil {
				http.Error(w, fmt.Sprintf("{\"error\":\"%v\"}", err), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "{\"status\":\"success\", \"id\":\"%s\", \"pid\":%d}", id, cmd.Process.Pid)

		})

		if err := http.ListenAndServe(":2388", mux); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Daemon shutting down")

}
