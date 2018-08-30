package manager

import (
	l "github.com/Yee2/Planet-Cup/ylog"
	"time"
)

var (
	UDPTimeout = 5 * time.Minute
	logger     = l.NewLogger("manager")
)
