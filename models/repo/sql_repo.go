package repo

import (
	"context"
	"database/sql"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"
)

type SqlRepository struct {
	config.Configuration
	logging.Logger
	Commands SqlCommands
	*sql.DB
	context.Context
}

type SqlCommands struct {
	// Init,
	// Seed,
	SelectUserBucketDict,
	SelectUserBuckets,
	SelectUserBucketKeyValue,
	GetPassage,
	GetBoltCatalog,
	SaveBoltData,
	CreateBucketIfNotExists *sql.Stmt
}
