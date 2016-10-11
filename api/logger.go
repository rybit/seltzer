package api

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/gommon/log"
)

// wrapper will adapt the logrus entries -> labstack's logger
type wrapper struct {
	*logrus.Entry
}

func (w wrapper) SetOutput(out io.Writer) {
	w.Logger.Out = out
}
func (w wrapper) SetLevel(l log.Lvl) {
	w.Logger.Level = logrus.Level(l)
}
func (w wrapper) Printj(json log.JSON) {
	w.Printf("%v", json)
}
func (w wrapper) Debugj(json log.JSON) {
	w.Debugf("%v", json)
}
func (w wrapper) Infoj(json log.JSON) {
	w.Infof("%v", json)
}
func (w wrapper) Warnj(json log.JSON) {
	w.Warnf("%v", json)
}
func (w wrapper) Errorj(json log.JSON) {
	w.Errorf("%v", json)
}
func (w wrapper) Fatalj(json log.JSON) {
	w.Fatalf("%v", json)
}
