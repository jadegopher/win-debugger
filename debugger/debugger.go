package debugger

import (
	"fmt"
	"log"
)

type Debugger interface {
	Debug() error
}

func (d *debugger) Debug() (err error) {
	if err = d.processRun(d.name); err != nil {
		return err
	}

	for {
		if err = d.waitForDebugEvent(); err != nil {
			log.Println(err)
			break
		}
		if d.continueDebugEvent() {
			continue
		}
		if err = d.threadContext(); err != nil {
			log.Println(err)
			break
		}
		fmt.Println(d.debugEvent, d.process)
	}

	return nil
}
