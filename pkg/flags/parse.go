package flags

import (
	"github.com/jessevdk/go-flags"
	"os"
)

func ParseFlags(config interface{}, args []string) error {
	_, err := flags.NewParser(config, flags.Default).ParseArgs(args)
	if err != nil {
		return err
	}
	return nil
}

func ParseOrFail(config interface{}, args []string) {
	err := ParseFlags(config, args)
	if err != nil {
		os.Exit(1)
	}
}
