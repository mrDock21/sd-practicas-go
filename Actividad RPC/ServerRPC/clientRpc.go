package ServerRPC

import (
	"encoding/json"
	"fmt"
	"net/rpc"
)

type Client struct {
	subjects *Subjects
	conn     *rpc.Client
}

func (c *Client) init() {
	c.subjects = &Subjects{
		Subjects: make(map[string]Students),
	}
	fmt.Println(c.subjects)
}

func (c *Client) Connect() bool {
	c.init()
	cl, err := rpc.Dial("tcp", "127.0.0.1"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return false
	}
	c.conn = cl
	return true
}

func (c *Client) Disconnect() {
	c.conn.Close()
	fmt.Println("---Disconnected---")
}

func (c *Client) ShowSubjects() {
	b, err := json.MarshalIndent(c.subjects, "", " ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}
}

func (c *Client) AddSubject(subject string) {
	args := Args{
		StrParams:     []string{subject},
		FloatParam:    0.0,
		SubjectsParam: *c.subjects,
	}
	reply := &Subjects{}
	err := c.conn.Call("Server.AddSubject", args, reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	*c.subjects = *reply
	fmt.Println("\n [Materia agregada!]")
}

func (c *Client) AddStudent(subject, student string, grade float64) {
	args := Args{
		StrParams:     []string{subject, student},
		FloatParam:    grade,
		SubjectsParam: *c.subjects,
	}
	reply := &Subjects{}
	err := c.conn.Call("Server.AddStudentToSubject", args, reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	*c.subjects = *reply
	fmt.Println("\n [Estudiante Agregado!]")
}

func (c *Client) ComputeSubjectGrade(subject string) {
	args := Args{
		StrParams:     []string{subject},
		FloatParam:    0.0,
		SubjectsParam: *c.subjects,
	}
	grade := 0.0
	reply := &grade
	err := c.conn.Call("Server.ComputeSubjectGrade", args, reply)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Materia: %s  [Promedio=%.2f]\n", subject, *reply)
	}
}

func (c *Client) ComputeStudentGrade(student string) {
	args := Args{
		StrParams:     []string{student},
		FloatParam:    0.0,
		SubjectsParam: *c.subjects,
	}
	reply := 0.0
	err := c.conn.Call("Server.ComputeStudentGrade", args, &reply)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Estudiante: %s  [Promedio=%.2f]\n", student, reply)
	}
}

func (c *Client) ComputeGeneralStudentGrade() {
	args := Args{
		StrParams:     []string{},
		FloatParam:    0.0,
		SubjectsParam: *c.subjects,
	}
	reply := 0.0
	err := c.conn.Call("Server.ComputeStudentsGeneralGrade", args, &reply)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("[Promedio-General=%.2f]\n", reply)
	}
}
