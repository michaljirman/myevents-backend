package dblayer

import (
	"github.com/michaljirman/myevents-backend/src/lib/persistence"
	"github.com/michaljirman/myevents-backend/src/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB  DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}
