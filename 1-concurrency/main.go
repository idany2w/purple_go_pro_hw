package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Генерирует случайное число в заданном диапазоне
func generateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// Первая горутина: генерирует случайные числа и передает их по каналу
func generateNumbers(wg *sync.WaitGroup, numbersChan chan<- int, count int) {
	defer wg.Done()
	defer close(numbersChan)

	for i := 0; i < count; i++ {
		number := generateRandomNumber(0, 100)
		numbersChan <- number
	}
}

// Вторая горутина: возводит числа в квадрат и передает результаты
func squareNumbers(wg *sync.WaitGroup, numbersChan <-chan int, resultsChan chan<- int) {
	defer wg.Done()
	defer close(resultsChan)

	for number := range numbersChan {
		square := number * number
		resultsChan <- square
	}
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(78.0 + time.Now().UnixNano())

	// Канал для передачи чисел от первой горутины ко второй
	numbersChan := make(chan int, 10)

	// Канал для передачи результатов от второй горутины в main
	resultsChan := make(chan int, 10)

	// WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup

	// Запускаем первую горутину
	wg.Add(1)
	go generateNumbers(&wg, numbersChan, 10)

	// Запускаем вторую горутину
	wg.Add(1)
	go squareNumbers(&wg, numbersChan, resultsChan)

	// Собираем результаты в main
	var results []int
	for result := range resultsChan {
		results = append(results, result)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Выводим все результаты
	for _, result := range results {
		fmt.Print(result, " ")
	}
}
