package connData

import "fmt"

type ConnData struct {
	MsgType int
	Task    Process
}

type Process struct {
	Id   uint
	Time uint64
}

func (p *Process) Update() {
	p.Time += 1
}

func (p *Process) ToString() string {
	return fmt.Sprintf("[ID: %d] Time: %d", p.Id, p.Time)
}
