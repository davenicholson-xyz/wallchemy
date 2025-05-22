package wallhaven

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/davenicholson-xyz/wallchemy/appcontext"
	"github.com/davenicholson-xyz/wallchemy/logger"
	"github.com/sirupsen/logrus"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func StartDaemon(app *appcontext.AppContext) {

	go func() {
		port := app.Config.GetIntWithDefault("port", 2388)

		logger.Log.WithField("port", port).Info("Started daemon")

		mux := http.NewServeMux()

		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			fmt.Fprintf(w, "{\"ping\":\"pong\"}")
		})

		mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			fmt.Fprintf(w, "{\"ping\":\"pong\"}")
		})

		mux.HandleFunc("GET /wp/{id}", func(w http.ResponseWriter, r *http.Request) {
			logger.Log.WithField("PATH", os.Getenv("PATH")).Info("Daemon environment PATH")
			enableCors(&w)

			id := r.PathValue("id")

			logger.Log.WithField("id", id).Debug("Daeomon request for wallpaper change")

			exePath, err := exec.LookPath("wallchemy")
			if err != nil {
				logger.Log.WithError(err).Error("wallchemy not found in PATH")
				http.Error(w, `{"error":"wallchemy not found in PATH"}`, http.StatusInternalServerError)
				return
			}
			logger.Log.WithField("path", exePath).Debug("Found wallchemy executable")

			output, err := exec.Command(exePath, "-id", id).CombinedOutput()
			if err != nil {
				logger.Log.WithFields(logrus.Fields{
					"error":  err,
					"output": string(output),
				}).Error("Failed to execute wallchemy command")
				http.Error(w, fmt.Sprintf("{\"error\":\"%v\"}", err), http.StatusInternalServerError)
				return
			}

			logger.Log.WithField("output", string(output)).Debug("Successfully changed wallpaper")

		})

		mux.HandleFunc("GET /kill", func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("Killing daemon")
			os.Exit(0)
		})

		portStr := fmt.Sprintf(":%d", port)

		if err := http.ListenAndServe(portStr, mux); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Daemon shutting down")
	log.Println("Daemon shutting down")

}
