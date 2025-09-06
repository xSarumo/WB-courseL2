package main

import "fmt"

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}

// Сначала выполнится функция test, и только выполнится defer на инкремент x. Так x это именованный параметр,
// то он сначала равен 1, а потом лямба функция прибавит к нему еще 1. Во втором же случаи мы возвращаем
// значение x, хотя defer выполниться он не окажет никакого влияния на результат. Поэтому вывод: 2 1
