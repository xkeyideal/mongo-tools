package options

// Struct encompassing all of the options that are reused across tools: "help",
// "version", verbosity settings, ssl settings, etc.
type ToolOptions struct {

	// The name of the tool
	AppName string

	Addrs []string

	// Sub-option types
	//*Connection
	Timeout             int
	TCPKeepAliveSeconds int

	//Auth
	Username string
	Password string
	Source   string
	//Mechanism string

	//Namespace
	// Specified database and collection
	DB         string
	Collection string

	// Force direct connection to the server and disable the
	// drivers automatic repl set discovery logic.
	Direct bool

	// ReplicaSetName, if specified, will prevent the obtained session from
	// communicating with any server which is not part of a replica set
	// with the given name. The default is to communicate with any server
	// specified or discovered via the servers contacted.
	ReplicaSetName string

	ReadTimeout int
	PoolLimit   int
}

// Ask for a new instance of tool options
func New(appName string) *ToolOptions {
	opts := &ToolOptions{
		AppName:     appName,
		PoolLimit:   2,
		ReadTimeout: 8,
		Direct:      false,
	}

	return opts
}

// Get the authentication database to use. Should be the value of
// --authenticationDatabase if it's provided, otherwise, the database that's
// specified in the tool's --db arg.
func (o *ToolOptions) GetAuthenticationDatabase() string {
	if o.Source != "" {
		return o.Source
	} else if o.DB != "" {
		return o.DB
	}
	return ""
}
