package golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println("hasil nil karena berbeda parent : ", contextF.Value("b"))
}

// Example Leak Goroutine
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()
	return destination
}
func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancelFunc := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		fmt.Println("Total Goroutine:", runtime.NumGoroutine())
		if n == 10 {
			break
		}
	}
	cancelFunc()
	time.Sleep(1 * time.Second)
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}

func CreateCounterContextTimeout(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) //simulasi slow response
			}
		}
	}()
	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancelFunc := context.WithTimeout(parent, 5*time.Second)
	defer cancelFunc()

	destination := CreateCounterContextTimeout(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		fmt.Println("Total Goroutine:", runtime.NumGoroutine())
		//trying looping never ending
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {

	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancelFunc := context.WithDeadline(parent, time.Now().Add(10*time.Second))
	defer cancelFunc()

	destination := CreateCounterContextTimeout(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		fmt.Println("Total Goroutine:", runtime.NumGoroutine())
		//trying looping never ending
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}
