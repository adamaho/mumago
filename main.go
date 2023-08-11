package main

import (
	"fmt"

	"mumago/internal/realtime"
)

func main() {
	message := realtime.Hello("world")
	fmt.Println(message)
}
