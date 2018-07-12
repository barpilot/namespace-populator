package controller

import (
	"time"
)

// Config is the controller configuration.
type Config struct {
	ResyncPeriod time.Duration
	Labels   string
	//Namespace    string
}
