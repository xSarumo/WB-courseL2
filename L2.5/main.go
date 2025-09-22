package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

// Программа выведет error потому что error это интерфейс и при
// err = test(), мы присваеваем значение интрефейсу в 2 поля в поле
// value = nil, а в поле type присваивается тип указатель на customError
// а интерейс == nil только при условии что оба поля nil
