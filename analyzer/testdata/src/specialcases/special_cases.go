package main

import "fmt"

type Count struct {
	Count1 int
	Count0 int
}

func main() {
	f := getAndIncrement()
	c := Count{
		Count0: f(),
		Count1: f(),
	}
	fmt.Printf("%+v\n", c)
}

func getAndIncrement() func() int {
	var count int
	return func() int {
		toReturn := count
		count++
		return toReturn
	}
}
