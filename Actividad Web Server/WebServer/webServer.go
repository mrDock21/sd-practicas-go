package WebServer

import (
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

func (s *Server) buildRoutes() {
	http.HandleFunc("/", s.root)
	http.HandleFunc("/add-subject", s.addSubject)
	http.HandleFunc("/add-grade", s.addStudent)
	//	to be able to retrieve css files
	staticHandler := http.FileServer(http.Dir("./views/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticHandler))
}

//	Root path
func (s *Server) root(res http.ResponseWriter, req *http.Request) {
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
	fmt.Println("[DEBUG]	Add subject")
	fmt.Println(s.mySubjects.Grades)
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"success.html"))
}

// Add student's grade
func (s *Server) addStudent(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		return
	}
	// add student to map
	subject := req.FormValue("subject-input2")
	_student := req.FormValue("student-input")
	grade, err := strconv.ParseFloat(req.FormValue("grade-input"), 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.mySubjects.Grades[subject][_student] = grade
	// respond with success
	fmt.Println("[DEBUG]	Add student")
	fmt.Println(s.mySubjects.Grades)
	respondWith(&res, HtmlRes, readFile(VIEWS_PATH+"success.html"))
}

func (s *Server) subjectGrade(res http.ResponseWriter, req *http.Request) {
	// only allow POST method
	if checkValidRequest(&res, req, "POST") == false {
		// not valid
		respondWith(&res, JsonRes, FAILURE_JSON)
		return
	}
	subject := req.FormValue("subject-input")
	grade := 0.0
	_studens, ok := s.mySubjects.Grades[subject]
	// check subject's existence
	if !ok {
		respondWith(&res, JsonRes, FAILURE_JSON+"{\"msg\":\"Subject doesn't exist!\"}")
		return
	}
	// compute grade
	for _, studentGrade := range _studens {
		grade += studentGrade
	}
	respondWith(&res, JsonRes, fmt.Sprintf("{\"Status\": 1,\"Data\":%.2f}", grade))
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
