package app

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/ranta0/rest-and-go/config"
	"github.com/ranta0/rest-and-go/domain/auth"
	"github.com/ranta0/rest-and-go/domain/role"
	"github.com/ranta0/rest-and-go/domain/user"
	"github.com/ranta0/rest-and-go/logging"
)

type App struct {
	Config    *config.Config
	DB        *gorm.DB
	Router    *chi.Mux
	Logger    *logging.LoggerFile
	LoggerAPI *logging.LoggerFile
}

func NewApp(cfg *config.Config) (*App, error) {
	r := chi.NewRouter()

	logAPP, err := logging.NewLoggerFile(cfg.LogDir, cfg.LogAPPFilename)
	if err != nil {
		return nil, err
	}

	logAPI, err := logging.NewLoggerFile(cfg.LogDir, cfg.LogAPIFilename)
	if err != nil {
		return nil, err
	}
	if cfg.Debug {
		logAPP.SetDebug()
		logAPI.SetDebug()
	}

	db, err := initDB(cfg, logAPP)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:    cfg,
		Router:    r,
		DB:        db,
		Logger:    logAPP,
		LoggerAPI: logAPI,
	}, nil
}

func initDB(cfg *config.Config, logger *logging.LoggerFile) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(
		&role.Role{},
		&user.User{},
		&auth.RevokedJWTToken{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
