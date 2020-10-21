package main

import (
	"fmt"
	"strconv"

	multimedia "./Multimedia"
)

func addImage(wc *multimedia.WebContent) {
	var title, format string

	fmt.Print("Ingrese el titulo de la image:")
	fmt.Scanln(&title)
	fmt.Print("Ingrese el formato de la image:")
	fmt.Scanln(&format)

	img := multimedia.Image{
		Title: title, Format: format,
	}
	img.Init()
	wc.Add(&img)
	fmt.Println("Se ha agregado la Imagen")
}

func addAudio(wc *multimedia.WebContent) {
	var title, format, durStr string
	var duration float32

	fmt.Print("Ingrese el titulo del audio:")
	fmt.Scanln(&title)
	fmt.Print("Ingrese el formato del audio:")
	fmt.Scanln(&format)
	fmt.Print("Ingrese la duración del audio:")
	fmt.Scanln(&durStr)

	f, _ := strconv.ParseFloat(durStr, 32)
	duration = float32(f)

	audio := multimedia.Audio{
		Title: title, Format: format, Duration: duration,
	}
	wc.Add(&audio)
	fmt.Println("Se ha agregado el Audio!\n")
}

func addVideo(wc *multimedia.WebContent) {
	var title, format, framesStr string
	var frames int32

	fmt.Print("Ingrese el titulo del video:")
	fmt.Scanln(&title)
	fmt.Print("Ingrese el formato del video:")
	fmt.Scanln(&format)
	fmt.Print("Ingrese los frames del video:")
	fmt.Scanln(&framesStr)

	i, _ := strconv.ParseInt(framesStr, 10, 32)
	frames = int32(i)

	video := multimedia.Video{
		Title: title, Format: format, Frames: frames,
	}
	wc.Add(&video)
	fmt.Println("Se ha agregado el Video!")
}

func main() {
	input := ""
	exit := false
	webContent := multimedia.WebContent{}

	for !exit {
		fmt.Println("+--------CONTENIDO-WEB--------+")
		fmt.Println("a) Agregar Image")
		fmt.Println("b) Agregar Audio")
		fmt.Println("c) Agregar Video")
		fmt.Println("d) Mostrar todo lo agregado")
		fmt.Println("e) Salir")
		fmt.Print("\nIngrese una opción:")
		fmt.Scanln(&input)

		switch {
		case input == "a":
			addImage(&webContent)
			break
		case input == "b":
			addAudio(&webContent)
			break
		case input == "c":
			addVideo(&webContent)
			break
		case input == "d":
			webContent.ShowAll()
			break
		default:
			fmt.Println("Saliendo...")
		}
		exit = input == "e"
	}
}
