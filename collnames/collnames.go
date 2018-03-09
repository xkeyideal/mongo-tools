package collnames

import (
	"github.com/xkeyideal/mongo-tools/common/db"
	"github.com/xkeyideal/mongo-tools/common/options"
)

type CollNames struct {
	Options *options.ToolOptions

	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewCollNames(opts *options.ToolOptions, sp *db.SessionProvider) *CollNames {
	return &CollNames{
		Options:         opts,
		SessionProvider: sp,
	}
}

func (cn *CollNames) Run() ([]string, error) {
	session, err := cn.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return cn.SessionProvider.CollectionNames(cn.Options.DB)
}
