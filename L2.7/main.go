package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

/* Вывод программы будет всегда разным из за особенности работы select.
Если в select есть 2 подходящих условия он выберет одно рандомное из них
Конструкция for select нужна для постоянной проверки не появилось ли что то новое
в каналах, что можно добавить. Таким образом алгоритм будет работать до тех пор пока
каналы не закроются.
Еще одним фактором случайного вывода чисел может послужить время, так как по сути
требуется время чтобы вызвать функцию и создать горутину. И в 2 одинаковых вызова функции,
дают всегда разное, но не координальное, время отработки -> какие то операции по добавлению
новых чисел в канал будут отрабатывать раньше.
*/
