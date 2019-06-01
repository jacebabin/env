package datastore

import "testing"

func TestDatastore_Option(t *testing.T) {
	logopt := InitLogDB()
	cacheopt := InitCacheDB()

	tests := []struct {
		name string
		opt option
	} {
		{"Log DB", logopt},
		{"Cache DB", cacheopt},
		{"Both DBs", nil}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewDatastore()
			if err != nil {
				t.Errorf("NewDatastore() error = %v", err)
				return
			}
			if tt.name == "Both DBs" {
				err = db.Option(logopt, cacheopt)
				if err != nil {
					t.Errorf("Option returned an err = %v", err)
				} else {
					err = ds.Option(tt.opt)
					if err != nil {
						t.Errorf("Option returned an err = %v", err)
					}
				}
			}
		})
	}
}