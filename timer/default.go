package timer

import (
	"github.com/panjf2000/ants/v2"
)

func InitDefaultTimer() error {
	Init(EDefaultTickMs, ants.Submit)
	GetTimingWheel().Run()
	return nil
}
