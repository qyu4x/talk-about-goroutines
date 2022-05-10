package goroutines

import (
	"fmt"
	"testing"
	"time"
)

func doStuff(till int, word string) {
	for i := 0; i < till; i ++ {
		fmt.Printf("%s:%d \n", word, i+1)
	}
}

func TestDoStuff(t *testing.T)  {
	go doStuff(10, "goroutines") // run asyncronus
	fmt.Println("done")

	go func(val string) {
		go doStuff(10, val )
	}("from anonymus func")

	time.Sleep(time.Second)

}

