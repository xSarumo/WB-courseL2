package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}

/* Программа выведет false, из за устройства интерфейса. Интерфейс по сути == nil при условие,
поле Value и Type == nil. Однако когда мы вызываем Foo() нам возвращается интерфейсный тип error,
для, которого выполняются все те же правила -> 1 из полей, а именно Type уже не будет == nil (Type == *os.PathError)
-> и весь интерфейс не равен nil
*/
