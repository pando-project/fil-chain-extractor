package lfx

import (
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
)

var log = logging.Logger("lfx")

var fxprinter fx.Printer = printer{}

type printer struct{}

func (printer) Printf(s string, args ...interface{}) {
	log.Debugf(s, args...)
}
