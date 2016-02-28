package settings

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type AppOptions struct {
	cfg   string
	gid   int
	nid   int
	roles string
}

func (o *AppOptions) Config() string {
	return o.cfg
}

func (o *AppOptions) Gid() int {
	return o.gid
}

func (o *AppOptions) Nid() int {
	return o.nid
}

func (o *AppOptions) Roles() []string {
	return strings.Split(o.roles, ",")
}

func (o *AppOptions) validate() []error {
	errors := make([]error, 0)
	if o.gid == 0 {
		errors = append(errors, fmt.Errorf("Gid can't be 0"))
	}

	if o.nid == 0 {
		errors = append(errors, fmt.Errorf("Nid can't be 0"))
	}

	return errors
}

var Options AppOptions

func init() {
	help := false
	flag.BoolVar(&help, "h", false, "Print this help screen")
	flag.StringVar(&Options.cfg, "c", "/etc/g8os/net.toml", "Path to config file")
	flag.IntVar(&Options.gid, "gid", 0, "Grid ID")
	flag.IntVar(&Options.nid, "nid", 0, "Node ID")
	flag.Parse()

	printHelp := func() {
		fmt.Println("core [options]")
		flag.PrintDefaults()
	}

	if help {
		printHelp()
		os.Exit(0)
	}

	if errors := Options.validate(); len(errors) != 0 {
		for _, err := range errors {
			fmt.Errorf("Validation Error: %s", err)
		}

		os.Exit(1)
	}
}