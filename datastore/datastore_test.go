package datastore

import "testing"

func TestNewDatastore(t *testing.T) {
	ds := new(Datastore)

	tests := []struct {
		name string
		want *Datastore
		wantErr bool
	} {
		{"Test 1", ds, false}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDatastore()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatastore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			appdb, err := got.DB(AppDB)
			if err != nil {
				t.Errorf("Error getting Database from Datastore = %v", err)
			}
			err = appdb.Ping() {
				if err != nil {
					t.Errorf("Error pinging database = %v", err)
				}
			}
		})
	}
}

func Test_NewDB(t *testing.T) {
	type args struct {
		n DBName
	}
	tests := []struct {
		name string
		args args
	} {
		{"App DB", args{AppDB}},
		{"Log DB", args{LogDB}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDB(tt.args.n)
			if err != nil {
				t.Errorf("Error from NewDB = %v", err)
			}
			err = db.Ping()
			if err !=  nil {
				t.Errorf("Error pinging database = %v", err)
			}
		})
	}
}