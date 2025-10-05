package main

import (
	"net/http"
	"record-services/internal/auth"
	"record-services/internal/config"
	"record-services/internal/middleware"
	"record-services/internal/migrations"
	"record-services/internal/repositories/user_repository"
	"record-services/pkg/database"
	"record-services/pkg/logger"
	"record-services/pkg/validator"
)

func main() {
	logger := logger.New()
	loggerApp := logger.Logger

	cfg := config.New()
	loggerApp.Info().Msg("Конфигурация успешно загружена")

	db, err := database.New(cfg.Db.Dsn)
	if err != nil {
		loggerApp.Fatal().Err(err).Msg("ошибка подключения к БД")
	}
	loggerApp.Info().Msg("Подключение к БД успешно")

	err = migrations.Migrate(db)
	if err != nil {
		loggerApp.Fatal().Err(err).Msg("ошибка применения миграций к БД")
	}
	loggerApp.Info().Msg("Миграции к БД успешно применены")

	// регистрация репозиториев
	userRepository := user_repository.NewUserRepository(db, loggerApp)

	//валидация
	validate := validator.NewValidate()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Главная страница"))
	})

	//регистрация routes
	authHandlers := auth.NewAuthHandlers(mux, loggerApp, userRepository, validate, cfg.Secret.HashSecret, cfg.Secret.JwtSecret)

	//middlewares
	middlewareAuth := middleware.AuthMiddleware(authHandlers)(mux)
	middlewareAuth = middleware.CORSMiddleware(middlewareAuth)

	// Установим уровень логирования из конфигурации
	logger.SetLogLevel(cfg.LogLevel)
	//server
	if err := http.ListenAndServe(cfg.Server.Listen, middlewareAuth); err != nil {
		loggerApp.Fatal().Err(err).Msg("ошибка запуска сервера")
	}

}
