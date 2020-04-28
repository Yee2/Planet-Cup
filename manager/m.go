package manager

import (
	l "hyacinth/ylog"
	"time"
)

var (
	UDPTimeout = 5 * time.Minute
	logger     = l.NewLogger("manager")
)
