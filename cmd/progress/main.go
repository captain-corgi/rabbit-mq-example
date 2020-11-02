package main

import (
	"fmt"
	"strings"
)

var (
	playlist = []string{
		"",
	}
	inprogress = "https://www.rabbitmq.com/tutorials/tutorial-six-go.html"
	finished   = []string{
		"",
		"https://www.rabbitmq.com/tutorials/tutorial-one-go.html",
		"https://www.rabbitmq.com/tutorials/tutorial-two-go.html",
		"https://www.rabbitmq.com/tutorials/tutorial-three-go.html",
		"https://www.rabbitmq.com/tutorials/tutorial-four-go.html",
		"https://www.rabbitmq.com/tutorials/tutorial-five-go.html",
	}
)

func main() {
	fmt.Printf("PLAYLIST:\t %+v\n", strings.Join(playlist, "\n\t"))
	fmt.Printf("WATCHING:\n \t%s\n", inprogress)
	fmt.Printf("FINISHED:\t %+v\n", strings.Join(finished, "\n\t"))
}
