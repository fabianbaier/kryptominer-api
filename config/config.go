package config

var (
	Version  = "UNSET"
	Revision = "UNSET"
)

type Config struct {
	DBURL		string
	Version     string
	Revision    string
	HomePath    string
	FlagVerbose bool
	FlagJSONLog bool
	FlagAPIPort int
}

func (c *Config) setDefaults() {
	c.DBURL = "127.0.0.1"
	c.Version = Version
	c.Revision = Revision
	c.FlagVerbose = true
	c.FlagJSONLog = false
	c.FlagAPIPort = 8787
}

func Configuration() (c Config) {
	c.setDefaults()
	return c
}