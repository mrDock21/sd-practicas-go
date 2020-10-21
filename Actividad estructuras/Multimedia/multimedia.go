package multimedia

import (
	"fmt"
	"math/rand"
)

type Multimedia interface {
	show()
}

type WebContent struct {
	Content []Multimedia
}

type Image struct {
	Title    string
	Format   string
	Channels struct {
		R []uint8
		G []uint8
		B []uint8
	}
}

type Audio struct {
	Title    string
	Format   string
	Duration float32
}

type Video struct {
	Title  string
	Format string
	Frames int32
}

func (wc *WebContent) Add(m Multimedia) {
	wc.Content = append(wc.Content, m)
}

func (wc *WebContent) ShowAll() {
	fmt.Println("+--CONTENIDO WEB--+")
	for _, c := range wc.Content {
		c.show()
		fmt.Println("")
	}
}

func (img *Image) Init() {
	sizes := []int{2, 8, 16, 32}
	size := sizes[rand.Intn(100)%len(sizes)]

	for i := 0; i < size; i++ {
		img.Channels.R = append(img.Channels.R, uint8(rand.Intn(255)))
		img.Channels.G = append(img.Channels.G, uint8(rand.Intn(255)))
		img.Channels.B = append(img.Channels.B, uint8(rand.Intn(255)))
	}
}

func (img *Image) show() {
	fmt.Println("+-----IMAGEN-------+")
	fmt.Printf("[TITULO]	%s\n", img.Title)
	fmt.Printf("[FORMATO]	%s\n", img.Format)
	fmt.Printf("[CANALES]	%s\n", img.Title+img.Format)
	fmt.Print("|	[R]		")
	fmt.Println(img.Channels.R)
	fmt.Print("|	[G]		")
	fmt.Println(img.Channels.G)
	fmt.Print("|	[B]		")
	fmt.Println(img.Channels.B)
}

func (video *Video) show() {
	fmt.Println("+-----VIDEO-------+")
	fmt.Printf("[TITULO]	%s\n", video.Title)
	fmt.Printf("[FORMATO]	%s\n", video.Format)
	fmt.Printf("[FRAMES]	%d\n", video.Frames)
}

func (audio *Audio) show() {
	fmt.Println("+-----AUDIO-------+")
	fmt.Printf("[TITULO]	%s\n", audio.Title)
	fmt.Printf("[FORMATO]	%s\n", audio.Format)
	fmt.Printf("[DURACION]	%.2fs\n", audio.Duration)
}
