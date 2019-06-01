package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/jacebabin/errors"
	_ "github.com/lib/pq" // pq driver calls for blank identifier
)

// DBName defines database name
type DBName int

const (
	// AppDB represents main application database
	AppDB DBName = iota

	// LogDB represents http logging database
	LogDB
)

// Environment variables for the App DB PostreSQL Database
const (
	envAppDBName     = "PG_APP_DBNAME"
	envAppDBUser     = "PG_APP_USERNAME"
	envAppDBPassword = "PG_APP_PASSWORD"
	envAppDBHost     = "PG_APP_HOST"
	envAppDBPort     = "PG_APP_PORT"
)

// Environment variables for the Log DB PostgreSQL Database
const (
	envLogDBName     = "PG_LOG_DBNAME"
	envLogDBUser     = "PG_LOG_USERNAME"
	envLogDBPassword = "PG_LOG_PASSWORD"
	envLogDBHost     = "PG_LOG_HOST"
	envLogDBPort     = "PG_LOG_PORT"
)

// Datastore struct stores common environment related items
type Datastore struct {
	appDB   *sql.DB
	logDB   *sql.DB
	cacheDB *redis.Pool
}

// NewDatastore initializes the datastore struct
func NewDatastore() (*Datastore, error) {
	const op errors.Op = "env/datastore/NewDatastore"

	// Get an AppDB (PostgreSQL)
	adb, err := NewDB(AppDB)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return &Datastore{appDB: adb, logDB: nil, cacheDB: nil}, nil
}

// BeginTx begins a *sql.Tx for the given db
func (ds Datastore) BeginTx(ctx context.Context, opts *sql.TxOptions, n DBName) (*sql.Tx, error) {
	const op errors.Op = "env/datastore/Datastore.BeginTx"

	switch n {
	case AppDB:
		// Calls the BeginTx method of the mainDB opened database
		mtx, err := ds.appDB.BeginTx(ctx, opts)
		if err != nil {
			return nil, errors.E(op, err)
		}
		return mtx, nil
	case LogDB:
		// Calls the BeginTx method of the mogDB opened database
		ltx, err := ds.logDB.BeginTx(ctx, opts)
		if err != nil {
			return nil, errors.E(op, err)
		}
		return ltx, nil
	default:
		return nil, errors.E(op, "Unexpected Database Name")
	}
}

// DB returns an initialized sql.DB given a database name
func (ds Datastore) DB(n DBName) (*sql.DB, error) {
	const op errors.Op = "env/datastore/Datastore.DB"

	switch n {
	case AppDB:
		if ds.appDB == nil {
			return nil, errors.E(op, "AppDB has not been initialized")
		}
		return ds.appDB, nil
	case LogDB:
		if ds.logDB == nil {
			return nil, errors.E(op, "LogDB has not been initialized")
		}
		return ds.logDB, nil
	default:
		return nil, errors.E(op, "Unexpected Database Name")
	}
}

// NewDB returns an open database handle of 0 or more underlying connections
func NewDB(n DBName) (*sql.DB, error) {
	const op errors.Op = "env/datastore/NewDB"

	// Get Datastore connection credentials from environment variables
	dbName := os.Getenv(dbEnvName(n))
	dbUser := os.Getenv(dbEnvUser(n))
	dbPassword := os.Getenv(dbEnvPassword(n))
	dbHost := os.Getenv(dbEnvHost(n))
	dbPort, err := strconv.Atoi(os.Getenv(dbEnvPort(n)))
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Craft connection string for database connection
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open the postgres database using the postgres driver (pq)
	// func Open(driverName, dataSourceName string) (*DB, error)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Call Ping to validate the newly opened database is actually alive
	if err = db.Ping(); err != nil {
		return nil, errors.E(op, err)
	}
	return db, nil
}

func dbEnvName(n DBName) string {
	switch n {
	case AppDB:
		return envAppDBName
	case LogDB:
		return envLogDBName
	default:
		return ""
	}
}

func dbEnvUser(n DBName) string {
	switch n {
	case AppDB:
		return envAppDBUser
	case LogDB:
		return envLogDBUser
	default:
		return ""
	}
}

func dbEnvPassword(n DBName) string {
	switch n {
	case AppDB:
		return envAppDBPassword
	case LogDB:
		return envLogDBPassword
	default:
		return ""
	}
}

func dbEnvHost(n DBName) string {
	switch n {
	case AppDB:
		return envAppDBHost
	case LogDB:
		return envLogDBHost
	default:
		return ""
	}
}

func dbEnvPort(n DBName) string {
	switch n {
	case AppDB:
		return envAppDBPort
	case LogDB:
		return envLogDBPort
	default:
		return ""
	}
}
