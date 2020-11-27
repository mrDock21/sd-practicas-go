package WebServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ResponseType int

const (
	CONN_IP_PORT = ":9000"
	VIEWS_PATH   = "./views/"
	FAILURE_JSON = "{\"Status\": 0,\"Data\":0}"
	// Codes used for errors
	HtmlRes ResponseType = iota
	JsonRes
)

type Server struct {
	mySubjects Subjects
}

type students map[string]float64

type Subjects struct {
	Grades map[string]students
}

func (s *Server) Serve() {

	s.mySubjects = Subjects{
		Grades: make(map[string]students),
	}

	s.buildRoutes()
	http.ListenAndServe(CONN_IP_PORT, nil)

	fmt.Println("[SERVER-RUNNING...]")
}

//	Sets routes' strings to functions
func (s *Server) buildRoutes() {
	http.HandleFunc("/", s.root)
	http.HandleFunc("/add-subject", s.addSubject)
	http.HandleFunc("/add-grade", s.addStudent)
	http.HandleFunc("/subject-grade", s.getSubjectGrade)
	http.HandleFunc("/subject-grade/compute", s.subjectGrade)
	http.HandleFunc("/student-grade", s.getStudentGrade)
	http.HandleFunc("/student-grade/compute", s.studentGrade)
	http.HandleFunc("/general-grade", s.getGeneralGrade)
	http.HandleFunc("/general-grade/compute", s.generalGrade)
	http.HandleFunc("/show-all", s.getShowAll)
	http.HandleFunc("/show-all/get", s.showAll)
	//	to be able to retrieve css files
	staticHandler := http.FileServer(http.Dir("./views/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticHandler))
}

//	Root path
func (s *Server) root(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[SERVER]	Request INDEX")
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"index.html"))
}

// Add subject
func (s *Server) addSubject(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		return
	}
	// add subject to map
	subject := req.FormValue("subject-input")
	s.mySubjects.Grades[subject] = make(students)
	// respond with success
	fmt.Println("[SERVER]	Request ADD-SUBJECT=" + subject)
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"success.html"))
}

// Add student's grade
func (s *Server) addStudent(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		return
	}
	// retrieve data
	subject := req.FormValue("subject-input2")
	_student := req.FormValue("student-input")
	grade, err := strconv.ParseFloat(req.FormValue("grade-input"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, ok := s.mySubjects.Grades[subject]; !ok {
		fmt.Println("[SERVER]	DENIED:Non-existant subject")
		respondException(&res, "No existe esa materia!")
		return
	}
	if _, ok := s.mySubjects.Grades[subject][_student]; ok {
		fmt.Println("[SERVER]	DENIED:Already-existant student")
		respondException(&res, "El alumno agregado ya existe!")
		return
	}
	s.mySubjects.Grades[subject][_student] = grade
	// respond with success
	fmt.Println("[SERVER]	Request ADD-STUDENT=" + _student + " TO=" + subject)
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"success.html"))
}

//	Responds with subject grade's html
func (s *Server) getSubjectGrade(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[SERVER]	Request SUBJECT-GRADE-HTML")
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"subject_grade.html"))
}

//	Computes subject's grade and responds with JSON
func (s *Server) subjectGrade(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	req.ParseMultipartForm(500)
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		respondWith(&res, JsonRes, FAILURE_JSON)
		return
	}
	subject := req.FormValue("subject-input")
	grade := 0.0
	_studens, ok := s.mySubjects.Grades[subject]
	fmt.Println("[SERVER]	Request SUBJECT-GRADE=" + subject)
	// check subject's existence
	if !ok || len(_studens) == 0 {
		respondWith(&res, JsonRes, FAILURE_JSON)
		return
	}
	// compute grade
	for _, studentGrade := range _studens {
		grade += studentGrade
	}
	respondWith(&res, JsonRes, fmt.Sprintf("{\"Status\": 1,\"Data\":%.2f}", grade/float64(len(_studens))))
}

//	Responds with student grade's html
func (s *Server) getStudentGrade(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[SERVER]	Request STUDENT-GRADE-HTML")
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"student_grade.html"))
}

//	Computes subject's grade and responds with JSON
func (s *Server) studentGrade(res http.ResponseWriter, req *http.Request) {
	// prepare form
	req.ParseMultipartForm(500)
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		respondWith(&res, JsonRes, FAILURE_JSON)
		return
	}
	student := req.FormValue("student-input")
	fmt.Println("[SERVER]	Request STUDENT-GRADE=" + student)
	// compute
	grade := s.computeGrade(student)
	respondWith(&res, JsonRes, fmt.Sprintf("{\"Status\": 1,\"Data\":%.2f}", grade))
}

//	Computes given student's grade
func (s *Server) computeGrade(student string) float64 {
	grade := 0.0
	subjCount := 0
	for _, studentGrades := range s.mySubjects.Grades {
		stndtGrade, ok := studentGrades[student]
		if ok {
			grade += stndtGrade
			subjCount++
		}
	}
	// avoid 0 division
	if subjCount == 0 {
		grade = 0.0
	} else {
		grade = grade / float64(subjCount)
	}
	return grade
}

//	Responds with general grade's html
func (s *Server) getGeneralGrade(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[SERVER]	Request GENERAL-GRADE-HTML")
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"general_grade.html"))
}

//	Computes general grade and responds with JSON
func (s *Server) generalGrade(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		respondWith(&res, JsonRes, FAILURE_JSON)
		return
	}
	fmt.Println("[SERVER]	Request GENERAL-GRADE")
	allStudents := s.getAllStudents()
	// compute
	grade := 0.0
	if len(allStudents) > 0 {
		for _, sGrade := range allStudents {
			grade += sGrade
		}
		grade = grade / float64(len(allStudents))
	}
	respondWith(&res, JsonRes, fmt.Sprintf("{\"Status\": 1,\"Data\":%.2f}", grade))
}

//	Returns all students with computed grades each
func (s *Server) getAllStudents() students {
	allStudents := make(students)
	for _, _students := range s.mySubjects.Grades {
		for student, _ := range _students {
			_, ok := allStudents[student]
			if !ok {
				// append non-existant students
				allStudents[student] = s.computeGrade(student)
			}
		}
	}
	return allStudents
}

//	Responds with ShowAll html
func (s *Server) getShowAll(res http.ResponseWriter, req *http.Request) {
	fmt.Println("[SERVER]	Request SHOW-ALL-HTML")
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"show_all.html"))
}

//	Responds with ShowAll html
func (s *Server) showAll(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		respondWith(&res, JsonRes, FAILURE_JSON)
		return
	}
	fmt.Println("[SERVER]	Request SHOW-ALL")
	respondWith(&res, JsonRes, s.getSubjectsAsJSON())
}

//	Returns string of subjects as JSON
func (s *Server) getSubjectsAsJSON() string {
	var _json string
	jsonBytes, err := json.Marshal(s.mySubjects.Grades)
	if err != nil {
		_json = FAILURE_JSON
		fmt.Println(err)
	} else {
		_json = string(jsonBytes)
	}
	return _json
}

//	Validates request's form and desired METHOD
//	Returns false if invalid
func checkValidRequest(res *http.ResponseWriter, req *http.Request, desiredMethod string) bool {
	valid := true
	if req != nil {
		if !isFormValid(req) {
			return !valid
		}
	}
	if desiredMethod != req.Method {
		http.NotFound(*res, req)
		return !valid
	}
	return valid
}

//	Validates request's form
//	Returns false if invalid
func isFormValid(req *http.Request) bool {
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//	Fills Http/JSON response with given data
func respondWith(res *http.ResponseWriter, rtype ResponseType, data string) {
	typeStr := ""

	switch rtype {
	case HtmlRes:
		typeStr = "text-html"
	case JsonRes:
		typeStr = "application/json"
	}

	(*res).Header().Set(
		"Content-type",
		typeStr,
	)
	fmt.Fprintf(*res, data)
}

//	Reads given file and returns its contents
func readFile(path string) string {
	text, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(text)
}

func respondException(res *http.ResponseWriter, msg string) {
	str := fmt.Sprintf(readFile(VIEWS_PATH+"exception.html"), msg)
	respondWith(res, HtmlRes, str)
}
