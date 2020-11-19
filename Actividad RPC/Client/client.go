package main

import (
	"fmt"
	"os"
	"os/exec"

	ServerRPC "../ServerRPC"
)

//	Clears console
func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

//	Add subject to inner map
func addSubject(c *ServerRPC.Client) {
	var input string
	fmt.Print("\n Nombre de la materia: ")
	fmt.Scanln(&input)
	c.AddSubject(input)
}

//	Add student's grade to subject
func addStudentSubject(c *ServerRPC.Client) {
	var subject, student string
	var grade float64
	fmt.Print("\n Nombre de la materia: ")
	fmt.Scanln(&subject)
	fmt.Print("\n Nombre del alumno: ")
	fmt.Scanln(&student)
	fmt.Print("\n Calificación: ")
	fmt.Scanln(&grade)
	c.AddStudent(subject, student, grade)
}

//	Computes subject's grade
func computeSubjectGrade(c *ServerRPC.Client) {
	var subject string
	fmt.Print("\n Nombre de la materia: ")
	fmt.Scanln(&subject)
	c.ComputeSubjectGrade(subject)
}

//	Computes student's grade
func computeStudentGrade(c *ServerRPC.Client) {
	var student string
	fmt.Print("\n Nombre del estudiante: ")
	fmt.Scanln(&student)
	c.ComputeStudentGrade(student)
}

//	Shows menu options
func menu(c *ServerRPC.Client) {
	var input string
	isActive := true

	for isActive {
		fmt.Println("+---------------MATERIAS---------------+")
		fmt.Println(" a) Agregar materia")
		fmt.Println(" b) Agregar calificación de alumno")
		fmt.Println(" c) Obtener promedio de alumno")
		fmt.Println(" d) Obtener promedio de todos los alumnos")
		fmt.Println(" e) Obtener promedio de materio")
		fmt.Println(" f) Ver todo")
		fmt.Println(" s) Salir")

		fmt.Print("\n Ingrese una opción: ")
		fmt.Scanln(&input)

		switch {
		case input == "a":
			addSubject(c)
		case input == "b":
			addStudentSubject(c)
		case input == "c":
			computeStudentGrade(c)
		case input == "d":
			c.ComputeGeneralStudentGrade()
		case input == "e":
			computeSubjectGrade(c)
		case input == "f":
			c.ShowSubjects()
		}
		isActive = input != "s"

		if isActive {
			fmt.Scanln(&input)
			clear()
		}
	}
	c.Disconnect()
}

func main() {
	client := ServerRPC.Client{}

	if client.Connect() {
		menu(&client)
	}
}
