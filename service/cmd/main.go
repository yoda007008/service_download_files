package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testtask/service/internal/handlers"
	"testtask/service/internal/service"
	"time"
)

func main() {
	manager := service.NewTaskManager("tasks.json", 3)
	if err := manager.Load(); err != nil {
		slog.Info("Ошибка загрузки состояния:", err)
	}
	go manager.Run()

	mux := http.NewServeMux()
	taskHandlers := handlers.NewTaskHandlers(manager)
	taskHandlers.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// todo graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Сервис запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Info("HTTP сервер:", err)
		}
	}()

	<-stop
	slog.Info("Остановка...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Info("Ошибка остановки сервера:", err)
	}

	manager.Save()
	manager.Stop()

	slog.Info("Сервис остановлен")
}
