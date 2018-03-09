package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/xkeyideal/mongo-tools/common/options"
	mgo "gopkg.in/mgo.v2"
)

type (
	sessionFlag uint32
)

// Session flags.
const (
	None      sessionFlag = 0
	Monotonic sessionFlag = 1 << iota
	DisableSocketTimeout
)

// Used to manage database sessions
type SessionProvider struct {

	// For connecting to the database
	connector DBConnector

	// used to avoid a race condition around creating the master session
	masterSessionLock sync.Mutex

	// the master session to use for connection pooling
	masterSession *mgo.Session

	// flags for generating the master session
	readPreference mgo.Mode
	readTimeout    time.Duration
	poolLimit      int
}

// Returns a session connected to the database server for which the
// session provider is configured.
func (self *SessionProvider) GetSession() (*mgo.Session, error) {
	self.masterSessionLock.Lock()
	defer self.masterSessionLock.Unlock()

	// The master session is initialized
	if self.masterSession != nil {
		return self.masterSession.Copy(), nil
	}

	// initialize the provider's master session
	var err error
	self.masterSession, err = self.connector.GetNewSession()
	if err != nil {
		return nil, fmt.Errorf("error connecting to db server: %v", err)
	}

	// update masterSession based on flags
	self.refresh()

	// copy the provider's master session, for connection pooling
	return self.masterSession.Copy(), nil
}

// Close closes the master session in the connection pool
func (self *SessionProvider) Close() {
	self.masterSessionLock.Lock()
	defer self.masterSessionLock.Unlock()
	if self.masterSession != nil {
		self.masterSession.Close()
	}
}

// SetReadPreference sets the read preference mode in the SessionProvider
// and eventually in the masterSession
func (self *SessionProvider) SetReadPreference(pref mgo.Mode) {
	self.masterSessionLock.Lock()
	defer self.masterSessionLock.Unlock()

	self.readPreference = pref

	if self.masterSession != nil {
		self.refresh()
	}
}

// refresh is a helper for modifying the session based on the
// session provider flags passed in with SetFlags.
// This helper assumes a lock is already taken.
func (self *SessionProvider) refresh() {
	// handle readPreference
	self.masterSession.SetMode(self.readPreference, true)
	self.masterSession.SetPoolLimit(self.poolLimit)
	self.masterSession.SetSyncTimeout(self.readTimeout)
}

// NewSessionProvider constructs a session provider but does not attempt to
// create the initial session.
func NewSessionProvider(opts *options.ToolOptions) (*SessionProvider, error) {
	// create the provider
	provider := &SessionProvider{
		readPreference: mgo.PrimaryPreferred,
		readTimeout:    time.Duration(opts.ReadTimeout) * time.Second,
		poolLimit:      opts.PoolLimit,
	}

	// create the connector for dialing the database
	provider.connector = getConnector(opts)

	// configure the connector
	err := provider.connector.Configure(opts)
	if err != nil {
		return nil, fmt.Errorf("error configuring the connector: %v", err)
	}
	return provider, nil
}

// Get the right type of connector, based on the options
func getConnector(opts *options.ToolOptions) DBConnector {
	return &VanillaDBConnector{}
}
