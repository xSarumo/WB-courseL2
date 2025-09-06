package main

import "fmt"

func main() {
	a := [5]int{76, 77, 78, 79, 80}
	var b []int = a[1:4] // Программа выведет значения с 1го индекса по 4й не включительно
	fmt.Println(b)
}
