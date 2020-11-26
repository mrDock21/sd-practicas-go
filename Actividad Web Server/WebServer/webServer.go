package WebServer

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	CONN_IP_PORT = ":9000"
	VIEWS_PATH   = "./views/"
)

type Server struct {
	mySubjects Subjects
}

type Subjects struct {
	Grades map[string]struct {
		Students map[string]float64
	}
}

func (s *Server) Serve() {
	s.buildRoutes()
	http.ListenAndServe(CONN_IP_PORT, nil)

	fmt.Println("[SERVER-RUNNING...]")
}

func (s *Server) buildRoutes() {
	http.HandleFunc("/", s.root)

	//	to be able to retrieve css files
	staticHandler := http.FileServer(http.Dir("./views/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", staticHandler))
}

//	Root path
func (s *Server) root(res http.ResponseWriter, req *http.Request) {
	respondWith(&res, readFile(VIEWS_PATH+"index.html"))
}

//	Fills Http response with given data
func respondWith(res *http.ResponseWriter, data string) {
	(*res).Header().Set(
		"Content-type",
		"text-html",
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
