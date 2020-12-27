package registry

import (
	"flowban/model"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	SchemaName = "NO_DB"
)

func (reg *AppRegistry) initializeDatabase() error {
	// DB Connection Configuration
	// Handles GORM

	db, err := gorm.Open(mysql.Open(dbDsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Error().Msg("Failed Connecting to database " + err.Error())
		return err
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	//logging info
	log.Info().Msg("Connected to Database with following configuration:" +
		"\n Database Dialect \t: " + dbDialect +
		"\n Database Host \t\t: " + dbHost + ":" + dbPort +
		"\n Database Name \t\t: " + dbDatabase)

	//auto migration from GORM
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.ScrumProject{},
		&model.ScrumProjectMember{},
		&model.SprintSession{},
		&model.SprintIssue{},
		&model.ScrumKanban{})

	if err != nil {
		return err
	}
	reg.dbConn = db
	return nil
}
