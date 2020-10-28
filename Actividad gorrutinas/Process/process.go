package processes

import (
	"fmt"
	"time"
)

type Process struct {
	Id   uint
	Time uint64
	done *bool
}

func (p *Process) Init() {
	isDone := false
	p.done = &isDone
}

func (p *Process) Run(isVisible chan bool) {
	isProcessShown := false
	for !*p.done {

		p.Time += 1

		if isProcessShown {
			fmt.Printf("[ID: %d] Time: %d\n", p.Id, p.Time)
		}

		select {
		case show := <-isVisible:
			isProcessShown = show
		default:
			// just keep running
		}
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Printf("[ID: %d] PROCESO TERMINADO", p.Id)
}

func (p *Process) Terminate() {
	if *p.done {
		fmt.Printf("[ID: %d] ESTE PROCESO YA ESTBA TERMINADO...", p.Id)
	}
	*p.done = true
}

func (p *Process) IsDone() bool {
	return *p.done
}
