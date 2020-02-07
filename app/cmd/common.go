package cmd

// Database struct describes cli flags and arguments that are necessary
// to open connection with database
type Database struct {
	User     string `long:"dbuser" env:"DBUSER" required:"false" description:"database access user"`
	Driver   string `long:"dbdriver" env:"DBDRIVER" required:"true" description:"database driver"`
	Password string `long:"dbpassword" env:"DBPASSWORD" required:"false" description:"database access password"`
	Source   string `long:"dbsource" env:"DBSOURCE" required:"true" description:"database source in format <dbname>@<host>:port"`
}

type Hashing struct {
	BcryptCost int `long:"bcryptcost" env:"BCRYPTCOST" required:"true" description:"number of hashing applied to string"`
}

// CommonCommander simplifies the delivery of common arguments and variables
type CommonCommander interface {
	SetCommonOptions(common CommonOptions)
	Execute(args []string) error
}

// CommonOptions for each command
type CommonOptions struct {
	AppName     string
	AppAuthor   string
	Version     string
	LoggerFlags int
}

// SetCommonOptions options for the command
func (c *CommonOptions) SetCommonOptions(common CommonOptions) {
	c.AppAuthor = common.AppAuthor
	c.AppName = common.AppName
	c.Version = common.Version
}
