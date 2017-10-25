package config

import "flag"

type Config struct {
	DBURL     string
	EthWallet string
	Verbose   bool
	Port      int
}

func Configuration() Config {

	c := Config{}
	flag.IntVar(&c.Port, "port", 8787, "port to listen on.")
	flag.StringVar(&c.DBURL, "db", "127.0.0.1:27017", "mongodb url to connect to.")
	flag.StringVar(&c.EthWallet, "ethwallet", "", "eth addr to scrape.")
	flag.BoolVar(&c.Verbose, "v", false, "verbose mode.")

	flag.Parse()
	return c
}
