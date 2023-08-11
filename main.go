package main

import (
	"fmt"

	"muma/internal/realtime"
)

func main() {
	message := realtime.Hello("world")
	fmt.Println(message)
}
