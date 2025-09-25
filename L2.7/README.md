# Что выведет программа?

Объяснить работу конвейера с использованием select.

package main

import (
  "fmt"
  "math/rand"
  "time"
)

func asChan(vs ...int) &lt;-chan int {
  c := make(chan int)
  go func() {
    for _, v := range vs {
      c &lt;- v
      time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
    }
  close(c)
}()
  return c
}

func merge(a, b &lt;-chan int) &lt;-chan int {
  c := make(chan int)
  go func() {
    for {
      select {
        case v, ok := &lt;-a:
          if ok {
            c &lt;- v
          } else {
            a = nil
          }
        case v, ok := &lt;-b:
          if ok {
            c &lt;- v
          } else {
            b = nil
          }
        }
        if a == nil &amp;&amp; b == nil {
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