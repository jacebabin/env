package datastore

import "github.com/jacebabin/errors"

// Options sets the options specified
func (ds *Datastore) Option(opts ...option) error {
	const op errors.Op = "env/datastore/Datastore.Option"
	for _, opt := range opts {
		err := opt(ds)
		if err != nil {
			return errors.E(op, err)
		}
	}
	return nil
}

type option func(*Datastore) error

// InitLogDB initializes a postgres database for logging HTTP requests
// To be used with github.com/jacebabin/httplog
func InitLogDB() option {
	const op errors.Op = "env/datastore/Datastore.InitLogDB"
	return func(ds *Datastore) error {
		// Get a LogDB (PostgreSQL)
		ldb, err := NewDB(LogDB)
		if err != nil {
			return errors.E(op, err)
		}
		ds.logDB = ldb
		return nil
	}
}

// InitCacheDB initializes the cache database (redis) for the
// Datastore
func InitCacheDB() option {
	return func(ds *Datastore) error {
		// Get a Redis Pool from redigo client
		cdb := NewCacheDB()
		ds.cacheDB = cdb
		return nil
	}
}
