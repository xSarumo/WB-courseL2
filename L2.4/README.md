## Что выведет программа?

Объяснить вывод программы.
``` go
func main() {
  ch := make(chan int)
  go func() {
    for i := 0; i &lt; 10; i++ {
    ch &lt;- i
  }
}()
  for n := range ch {
    println(n)
  }
}
```