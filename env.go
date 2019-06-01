package env

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/jacebabin/env/datastore"
	"github.com/jacebabin/errors"
	"github.com/rs/zerolog"
)

// Name is the environment Name int representation
// Using iota, 1 (Production) is the lowest,
// 2 (Staging) is 2nd lowest, and so on...
type Name uint8

// Name of environment.
const (
	Production Name = iota + 1 // Production (1)
	Staging                    // Staging (2)
	QA                         // QA (3)
	Dev                        // Dev (4)
)

func (n Name) String() string {
	switch n {
	case Production:
		return "Production"
	case Staging:
		return "Staging"
	case QA:
		return "QA"
	case Dev:
		return "Dev"
	}
	return "unknown_name"
}

// Env struct stores common environment related items
type Env struct {
	// Environment Name (e.g. Production, QA, etc.)
	Name Name
	// multiplex router
	Router *mux.Router
	// Datastore struct containing AppDB (PostgreSQL)
	// LogDb (PostgreSQL) and CacheDB (Redis)
	DS *datastore.Datastore
	// Logger
	Logger zerolog.Logger
}

// NewEnv initialized the Env struct
func NewEnv(name Name, lvl zerolog.Level) (*Env, error) {
	const op errors.Op = "env/NewEnv"

	if name.String() == "unknown_name" {
		return nil, errors.E(op, "Unknown env.Name input")
	}

	if lvl.String() == "" {
		return nil, errors.E(op, "Unknown logger level input")
	}

	// setup logger
	log := newLogger(lvl)

	// open db connection pools
	dstore, err := datastore.NewDatastore()
	if err != nil {
		return nil, errors.E(op, err)
	}

	// create a new mux (multiplex) router
	rtr := mux.NewRouter()

	env := &Env{Name: name, Router: rtr, DS: dstore, Logger: log}

	return env, nil
}

// newSubrouter adds any subRouters that you'd like to have as part of
// every request, i.e. I always want to be sure that every request has
// "/api" as part of it's path prefix without having to put it into
// every handle path in my various routing functions
func newSubrouter(rtr *mux.Router) *mux.Router {
	sRtr := rtr.PathPrefix("/api").Subrouter()
	return sRtr
}

// newLogger sets up the zerologger.Logger
func newLogger(lvl zerolog.Level) zerolog.Logger {
	// empty string for TimeFieldFormat will write logs with UNIX time
	zerolog.TimeFieldFormat = ""
	// set logging level based on input
	zerolog.SetGlobalLevel(lvl)
	// start a new logger with Stdout as the target
	lgr := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return lgr
}
