package log

import (
	"log"
	"os"
)

var (
	Default = log.New(os.Stdout, "[karma-chameleon] ", log.LstdFlags)
	Error   = log.New(os.Stderr, "[karma-chameleon] ", log.LstdFlags)
	Slack   = log.New(os.Stdout, "[slack] ", log.LstdFlags)
)
