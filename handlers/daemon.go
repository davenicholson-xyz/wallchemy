package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davenicholson-xyz/wallchemy/app"
	"github.com/sirupsen/logrus"
)

type DaemonServer struct {
	server *http.Server
	appCtx *app.AppContext
}

func NewDaemonServer(appCtx *app.AppContext) *DaemonServer {
	mux := http.NewServeMux()

	ds := &DaemonServer{
		appCtx: appCtx,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", appCtx.Config.Port),
			Handler: mux,
		},
	}

	ds.setupRoutes(mux)

	return ds
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func (ds *DaemonServer) setupRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/wp/{id}", ds.handleWallpaper)
	mux.HandleFunc("/health", ds.handleHealth)
	mux.HandleFunc("/kill", ds.handleKill)

}

func (ds *DaemonServer) handleWallpaper(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	app.Logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"remote": r.RemoteAddr,
	}).Debug("/wp accessed")

	id := r.PathValue("id")

	ds.appCtx.Config.ID = id

	output, err := handleSingle(ds.appCtx)
	if err != nil {
		msg := err.Error()
		if err.Error() == "401 - Unauthorized" {
			msg = "API key required to access this wallpaper\n"
		}
		errMsg := map[string]string{"error": msg}
		jsonErr, _ := json.Marshal(errMsg)
		http.Error(w, string(jsonErr), http.StatusInternalServerError)
	} else {
		msg := map[string]string{"success": output}
		jsonMsg, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMsg)
	}

}

func (ds *DaemonServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	hasAPIkey := "false"
	apikey := ds.appCtx.Config.APIKey
	if apikey != "" {
		hasAPIkey = "true"
	}

	msg := map[string]string{"apikey": string(hasAPIkey)}
	jsonMsg, _ := json.Marshal(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMsg)
}

func (ds *DaemonServer) handleKill(w http.ResponseWriter, r *http.Request) {
	app.Logger.Info("Stopping daemon")
	os.Exit(0)
}

func (ds *DaemonServer) Start() error {
	// Start server in a goroutine
	go func() {
		app.Logger.WithFields(logrus.Fields{
			"port": ds.appCtx.Config.Port,
			"addr": ds.server.Addr,
		}).Info("Starting daemon server")

		if err := ds.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.WithError(err).Fatal("Server failed to start")
		}
	}()

	app.Logger.WithField("port", ds.appCtx.Config.Port).Info("Wallchemy daemon started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	app.Logger.WithField("signal", sig.String()).Info("Received shutdown signal")

	return ds.Shutdown()
}

func (ds *DaemonServer) Shutdown() error {
	app.Logger.Info("Shutting down daemon server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := ds.server.Shutdown(ctx); err != nil {
		app.Logger.WithError(err).Error("Server forced to shutdown")
		return err
	}

	app.Logger.Info("Daemon server stopped gracefully")
	return nil
}

func HandleDaemon(appCtx *app.AppContext) (string, error) {
	if appCtx.Config.Daemon == "" {
		return "", fmt.Errorf("daemon flag is empty")
	}

	switch appCtx.Config.Daemon {
	case "start":
		return LaunchDaemon(appCtx)
	case "stop":
		return stopDaemon(appCtx)
	case "status":
		return statusDaemon(appCtx)
	default:
		return "", fmt.Errorf("unknown daemon command: %s (use: start, stop, status)", appCtx.Config.Daemon)
	}
}

func startDaemon(appCtx *app.AppContext) (string, error) {
	if appCtx.Config.APIKey != "" {
		appCtx.WallhavenAPI.APIKey(appCtx.Config.APIKey)
	}

	server := NewDaemonServer(appCtx)

	if err := server.Start(); err != nil {
		return "", fmt.Errorf("failed to start daemon: %w", err)
	}

	return fmt.Sprintf("Daemon stopped"), nil
}

func stopDaemon(appCtx *app.AppContext) (string, error) {
	killStr := fmt.Sprintf("http://localhost:%d/kill", appCtx.Config.Port)

	_, err := http.Get(killStr)
	if err != nil {
		app.Logger.Debug("Could not kill daemon with GET request")
	}

	app.Logger.Info("Daemon process stopped")

	return "Wallchemy daemon stopped", nil
}

func statusDaemon(appCtx *app.AppContext) (string, error) {
	healthURL := fmt.Sprintf("http://localhost:%d/health", appCtx.Config.Port)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(healthURL)
	if err != nil {
		return "Daemon is not running", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return fmt.Sprintf("Daemon is running on port %d", appCtx.Config.Port), nil
	}
	return "Daemon status unknown", nil
}
