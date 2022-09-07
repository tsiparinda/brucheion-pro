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

	driver := config.GetStringDefault("sql:driverName", "postgres")
	connectionUrl, found := config.GetString("sql:connectionUrl")

	var connectionStr string
	var err error
	if !found {
		logger.Panic("openDB: Cannot read SQL connection string from config")
	} else {
		logger.Debug("openDB: found SQL connection string from config")
	}
	connectionStr, err = pq.ParseURL(connectionUrl)
	if err != nil {
		logger.Panic("openDB: Error converting SQL URL connection from config to connection string")
	}
	logger.Infof("openDB: Connection string: ", connectionStr)

	if db, err = sql.Open(driver, connectionStr); err == nil {
		logger.Debugf("openDB: db: ", db)

		loadMigrations(config, logger)

		commands = loadCommands(db, config, logger)
		logger.Debug("openDB: SQL commands loaded")

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
		logger.Debugf("loadCommands: Loading SQL command: %v", commandName)
		stmt := prepareCommand(db, commandName, config, logger)
		commandVal.Field(i).Set(reflect.ValueOf(stmt))
	}
	return commands
}

func prepareCommand(db *sql.DB, command string, config config.Configuration,
	logger logging.Logger) *sql.Stmt {
	filename, found := config.GetString("sql:commands:" + command)
	if !found {
		logger.Panicf("prepareCommand: Config does not contain location for SQL command: %v",
			command)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Panicf("prepareCommand: Cannot read SQL command file: %v", filename)
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

	logger.Debugf("loadMigrations: begin...")
	migrations_path := config.GetStringDefault("sql:migrationsPath", "file://./sql/migrations")

	connectionUrl, _ := config.GetString("sql:connectionUrl")
	logger.Debugf("loadMigrations: migrate input: ", connectionUrl, migrations_path)
	if m, err := migrate.New(migrations_path, connectionUrl); err == nil {

		logger.Debugf("loadMigrations: migrating: ", m, err)

		if config.GetBoolDefault("sql:alwaysReset", true) {
			logger.Debugf("loadMigrations: alwaysReset is true, downing migrate: ", m, err)
			if config.GetBoolDefault("sql:migrationsForce", false) {
				version := config.GetIntDefault("sql:migrationsVersion", -1)
				if err := m.Force(version); err != nil {
					logger.Debugf("loadMigrations: error in Force version %v: ", version, err)
				}
			}
			if err := m.Down(); err != nil {
				logger.Debugf("loadMigrations: downing migrate ends with error: ", err)
			}

		}
		logger.Debugf("loadMigrations: start to up migrations...")
		//  m.Force(1) if you want to explicitly set the number of migrations to run and reset dirty flag
		// m.Steps(1) for up to one step
		version, dirty, err := m.Version()
		logger.Debugf("loadMigrations: before run migration: version %d, dirty %v, error %s", version, dirty, err)
		if err := m.Up(); err != nil {
			logger.Debugf("loadMigrations: up migrations ends with error: ", err)
		}
		version, dirty, err = m.Version()
		logger.Debugf("loadMigrations: after run migration: version %d, dirty %v, error %s", version, dirty, err)
	} else {
		logger.Debugf("loadMigrations: Error migrate:  ", err)
	}
	//return err

	return
}

// for use for start from service
// //services.Call(func(repo models.Repository) { repo.LoadMigrations() })
// doesn't work as needed, because starts after load sql-queryes
func (repo *SqlRepository) LoadMigrations() {
	repo.Logger.Debugf("LoadMigrations: begin...")
	migrations_path := repo.Configuration.GetStringDefault("sql:migrationsPath", "file://./sql/migrations")

	connectionUrl, _ := repo.Configuration.GetString("sql:connectionUrl")
	repo.Logger.Debugf("LoadMigrations: migrate input: ", connectionUrl, migrations_path)
	if m, err := migrate.New(migrations_path, connectionUrl); err == nil {

		repo.Logger.Debugf("LoadMigrations: migrate: ", m, err)

		if repo.Configuration.GetBoolDefault("sql:alwaysReset", true) {
			m.Down()
		}

		m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	} else {
		repo.Logger.Debugf("LoadMigrations: Error migrate:  ", err)
	}
	//return err

	return
}
