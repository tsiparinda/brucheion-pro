package repo

import (
	"database/sql"
	"os"
	"reflect"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"
)

func openDB(config config.Configuration, logger logging.Logger) (db *sql.DB,
	commands *SqlCommands) {

	driver := config.GetStringDefault("sql:driver_name", "postgres")
	connectionUrl, found := config.GetString("sql:connection_url")

	var connectionStr string
	var err error
	if !found {
		logger.Panic("Cannot read SQL connection string from config")
	} else {
		logger.Debug("openDB.found SQL connection string from config")
	}
	connectionStr, err = pq.ParseURL(connectionUrl)
	if err != nil {
		logger.Panic("Error converting SQL URL connection from config to connection string")
	}
	logger.Infof("Connection string: ", connectionStr)

	if db, err = sql.Open(driver, connectionStr); err == nil {
		logger.Debugf("openDB.db: ", db)
		commands = loadCommands(db, config, logger)
		logger.Debug("openDB.SQL commands loaded")

	} else {
		logger.Panic(err.Error())
	}
	return
}

func loadCommands(db *sql.DB, config config.Configuration,
	logger logging.Logger) (commands *SqlCommands) {
	commands = &SqlCommands{}
	commandVal := reflect.ValueOf(commands).Elem()
	commandType := reflect.TypeOf(commands).Elem()
	for i := 0; i < commandType.NumField(); i++ {
		commandName := commandType.Field(i).Name
		logger.Debugf("Loading SQL command: %v", commandName)
		stmt := prepareCommand(db, commandName, config, logger)
		commandVal.Field(i).Set(reflect.ValueOf(stmt))
	}
	return commands
}

func prepareCommand(db *sql.DB, command string, config config.Configuration,
	logger logging.Logger) *sql.Stmt {
	filename, found := config.GetString("sql:commands:" + command)
	if !found {
		logger.Panicf("Config does not contain location for SQL command: %v",
			command)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Panicf("Cannot read SQL command file: %v", filename)
	} else {
		logger.Debug("prepareCommand: sql file readed")
	}
	statement, err := db.Prepare(string(data))
	if err != nil {
		logger.Panicf(err.Error())
	}
	return statement
}

// run in RegisterSqlRepositoryService
func loadMigrations(config config.Configuration, logger logging.Logger) {

	migrations_path := config.GetStringDefault("sql:migrations_path", "file://./sql/migrations")

	connectionUrl, _ := config.GetString("sql:connection_url")
	logger.Debugf("migrate input: ", connectionUrl, migrations_path)
	if m, err := migrate.New(migrations_path, connectionUrl); err == nil {

		logger.Debugf("migrate: ", m, err)

		if config.GetBoolDefault("sql:always_reset", true) {
			m.Down()
		}

		m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	} else {
		logger.Debugf("Error migrate:  ", err)
	}
	//return err

	return
}

// for use for start from service
// //services.Call(func(repo models.Repository) { repo.LoadMigrations() })
// doesn't work as needed, because starts after load sql-queryes
func (repo *SqlRepository) LoadMigrations() {

	migrations_path := repo.Configuration.GetStringDefault("sql:migrations_path", "file://./sql/migrations")

	connectionUrl, _ := repo.Configuration.GetString("sql:connection_url")
	repo.Logger.Debugf("migrate input: ", connectionUrl, migrations_path)
	if m, err := migrate.New(migrations_path, connectionUrl); err == nil {

		repo.Logger.Debugf("migrate: ", m, err)

		if repo.Configuration.GetBoolDefault("sql:always_reset", true) {
			m.Down()
		}

		m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	} else {
		repo.Logger.Debugf("Error migrate:  ", err)
	}
	//return err

	return
}
