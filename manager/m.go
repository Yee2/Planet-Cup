package manager

import (
	"time"
	l "github.com/Yee2/Planet-Cup/ylog"
)
var (
	UDPTimeout = 5*time.Minute
	logger = l.NewLogger("manager")
)