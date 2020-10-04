package main

import "fmt"

type greeting string

func (g greeting) Greet() {
	fmt.Println("Hello Universe")
}

// exported as symbol named "Greeter"
// nolint
var Greeter greeting
