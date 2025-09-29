package main

import (
	"record-services/internal/config"
	"record-services/internal/migrations"
	"record-services/pkg/database"
	"record-services/pkg/logger"
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

	// Установим уровень логирования из конфигурации
	logger.SetLogLevel(cfg.LogLevel)

}
