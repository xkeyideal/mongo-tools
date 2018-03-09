package collindexes

import (
	"mongo-tools/common/db"
	"mongo-tools/common/options"

	"gopkg.in/mgo.v2"
)

type CollIndexes struct {
	Options *options.ToolOptions

	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewCollIndexes(opts *options.ToolOptions, sp *db.SessionProvider) *CollIndexes {
	return &CollIndexes{
		Options:         opts,
		SessionProvider: sp,
	}
}

func (ci *CollIndexes) Indexes() ([]mgo.Index, error) {
	return ci.SessionProvider.Indexes(ci.Options.DB, ci.Options.Collection)
}

func (ci *CollIndexes) DropIndexName(name string) error {
	return ci.SessionProvider.DropIndexName(ci.Options.DB, ci.Options.Collection, name)
}
