package main

import (
	"fmt"
	"os"
	"os/exec"

	process "./Process"
)

func clear() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}

func ProcessBuilder() func() process.Process {
	count := 0
	return func() process.Process {
		count += 1
		process := process.Process{Id: uint(count), Time: 0}
		process.Init()
		return process
	}
}

func showAllProcesses(showAll chan bool, s []process.Process) {
	isShowingSomething := false
	// update channel for all processes
	for _, v := range s {
		if !v.IsDone() {
			isShowingSomething = true
			showAll <- true
		}
	}

	if !isShowingSomething {
		fmt.Print("NO HAY NINGÚN PROCESO EN EJECUCIÓN...")
	}
	// wait user input...
	fmt.Scanln()

	if isShowingSomething {
		// update channel for all processes
		for _, v := range s {
			if !v.IsDone() {
				showAll <- false
			}
		}
	}
}

func main() {
	var input string
	var processes []process.Process

	nextProcess := ProcessBuilder()
	showAll := make(chan bool)

	showMenu := true

	for showMenu {
		fmt.Println("+---------ADMINISTRADOR PROCESOS---------+")
		fmt.Println("|	a) Agregar nuevo proceso")
		fmt.Println("|	b) Mostrar procesos")
		fmt.Println("|	c) Terminar proceso")
		fmt.Println("|	d) Salir")
		fmt.Print("\n		Ingrese una opción:")
		fmt.Scanln(&input)

		switch {
		case input == "a":
			process := nextProcess()
			processes = append(processes, process)
			// start goroutine
			go process.Run(showAll)
			fmt.Printf("	Se agregó el proceso %d", process.Id)
			fmt.Scanln()
		case input == "b":
			showAllProcesses(showAll, processes)
		case input == "c":
			var index int
			fmt.Print("Ingrese el ID del proceso:")
			fmt.Scanln(&index)
			processes[index-1].Terminate()
			fmt.Scanln()
		case input == "d":
			fmt.Println("	Saliendo...")
		}
		showMenu = input != "d"
		if showMenu {
			clear()
		}
	}
}
