package app

type Log struct {
	Level  string `long:"level" description:"Log level" choice:"DEBUG" choice:"INFO" default:"DEBUG"`
	Format string `long:"format" description:"Format of message logs" choice:"json" choice:"text" default:"json"`
}

type Config struct {
	Address string `long:"address" description:"Http address" default:"0.0.0.0:8080"`
	Log     Log    `group:"log" namespace:"log"`
	Config  string `long:"config" required:"true" description:"Path to scraper scrape file"`
}
