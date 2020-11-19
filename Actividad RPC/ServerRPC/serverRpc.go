package ServerRPC

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

const (
	CONN_PORT = ":9999"
)

type Server struct{}

type Args struct {
	StrParams     []string
	FloatParam    float64
	SubjectsParam Subjects
}

type Subjects struct {
	//	key := subject Name, value := students (in subject)
	Subjects map[string]Students
}

type Students struct {
	//	key := sutdent Name, value := grade
	Grades map[string]float64
}

//	Starts server
func (s *Server) Serve() {
	rpc.Register(new(Server))
	rpcServer, err := net.Listen("tcp", CONN_PORT)

	if err != nil {
		fmt.Println(err)
		return
	}
	// start listening
	fmt.Println("[DEBUG]	Server listening...")
	for {
		conn, err := rpcServer.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

//	Computes given student's grade
func (s *Server) computeStudentGrade(name string, subjects map[string]Students, res *float64) error {
	totalSubjects, total := 0, 0.0
	// compute
	for _, students := range subjects {
		grade, ok := students.Grades[name]
		if ok {
			total += grade
			totalSubjects += 1.0
		}
	}
	// check 0 division
	if totalSubjects == 0 {
		return errors.New("Student" + name + " is non-existent in all subjects!")
	}
	*res = float64(total / float64(totalSubjects))
	return nil
}

//	Adds a subject to given map in Args
func (s *Server) AddSubject(args Args, reply *Subjects) error {
	fmt.Println("\n-----CALL ADD SUBJECT------")
	if len(args.StrParams) <= 0 {
		return errors.New("Needed subject name for AddSubject OP")
	}

	args.SubjectsParam.Subjects[args.StrParams[0]] = Students{
		Grades: make(map[string]float64),
	}

	*reply = args.SubjectsParam
	return nil
}

//	Adds given student's grade to given subject
func (s *Server) AddStudentToSubject(args Args, reply *Subjects) error {
	fmt.Println("\n-----CALL ADD STUDENT TO SUBJECT------")
	// args should have something
	if len(args.StrParams) <= 0 {
		return errors.New("Needed subject & student names for AddStudent OP")
	}
	if len(args.SubjectsParam.Subjects) == 0 {
		return errors.New("No Subjects have been added! For AddStudent OP")
	}
	// retrieve names
	subjectName, studentName := args.StrParams[0], args.StrParams[1]
	// check first is subject exists
	subject, isSubjectOK := args.SubjectsParam.Subjects[subjectName]

	if isSubjectOK {
		// now check student doesn't exists (because we are adding a new one)
		_, studentExists := subject.Grades[studentName]

		if !studentExists {
			// we add it
			args.SubjectsParam.Subjects[subjectName].Grades[studentName] = args.FloatParam
		} else {
			return errors.New("Student already exists!")
		}
	} else {
		return errors.New("Subject doesn't exists!")
	}
	*reply = args.SubjectsParam
	return nil
}

//	Computes a single subject's grade
func (s *Server) ComputeSubjectGrade(args Args, reply *float64) error {
	fmt.Println("\n-----CALL SUBJECT GRADE------")
	if len(args.StrParams) <= 0 {
		return errors.New("No subject name passed for ComputeSubjectGrade OP")
	}
	if args.SubjectsParam.Subjects == nil {
		return errors.New("No Subjects have been added! For ComputeSubjectGrade OP")
	}
	// check this subject
	students, subjectExists := args.SubjectsParam.Subjects[args.StrParams[0]]
	if !subjectExists {
		return errors.New("No subject exists with given name=" + args.StrParams[0])
	}
	if len(students.Grades) == 0 {
		return errors.New("No students in given subject name=" + args.StrParams[0])
	}
	// compute average
	total := 0.0
	for _, grade := range students.Grades {
		total += grade
	}

	*reply = float64(total / float64(len(students.Grades)))
	return nil
}

//	Computes a single student's grade
func (s *Server) ComputeStudentGrade(args Args, reply *float64) error {
	fmt.Println("\n-----CALL STUDENT GRADE------")
	if len(args.StrParams) <= 0 {
		return errors.New("No student name passed for StudentGrade OP")
	}
	if args.SubjectsParam.Subjects == nil {
		return errors.New("No Subjects have been added! For StudentGrade OP")
	}
	// retrieve
	subjects := args.SubjectsParam.Subjects
	studentName := args.StrParams[0]
	// compute
	return s.computeStudentGrade(studentName, subjects, reply)
}

//	Computes general grade of all students (of all subjects combined)
func (s *Server) ComputeStudentsGeneralGrade(args Args, reply *float64) error {
	fmt.Println("\n-----CALL STUDENT GENERAL GRADE------")
	if len(args.SubjectsParam.Subjects) == 0 {
		return errors.New("No data passed to process for StudentsGeneralGrade OP")
	}
	if args.SubjectsParam.Subjects == nil {
		return errors.New("No Subjects have been added! For StudentsGeneralGrade OP")
	}
	total, aux := 0.0, 0.0
	resAux := &aux
	students := make(map[string]float64)
	subjects := args.SubjectsParam.Subjects

	// fill students with all subject's students
	for _, subjectStudents := range subjects {
		for studentName, _ := range subjectStudents.Grades {
			_, ok := students[studentName]
			if !ok {
				students[studentName] = 0.0
			}
		}
	}
	// compute grade for each student found
	for studentName, _ := range students {
		err := s.computeStudentGrade(studentName, subjects, resAux)
		if err != nil {
			return err
		}
		total += *resAux
	}
	*reply = total / float64(len(students))
	return nil
}
