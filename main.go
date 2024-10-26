package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type NumberArray struct {
	Numbers []int `json:"numbers"`
}

func main() {

	logFile, err := os.Create("app.log")
	defer logFile.Close()

	infoLog := log.New(io.MultiWriter(os.Stdout, logFile), "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(io.MultiWriter(os.Stdout, logFile), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	source := flag.String("source", "file", "Укажите источник данных (файл или stdin)")
	flag.Parse()

	if *source == "" {
		errorLog.Fatalf("Пожалуйста, укажите либо 'file', либо 'stdin'")
	}

	// Задание 2.1  Чтение из файла JSON с массивом чисел.
	var numbers *NumberArray
	switch *source {
	case "file":
		{
			numbers, err = readJSONFromFile("numbers.json")
			if err != nil {
				errorLog.Fatalf("Шаг 1: Чтение из файла JSON с массивом чисел. %v", err)
			}
			infoLog.Printf("Шаг 1: Чтение из файла JSON с массивом чисел. Чисел: %d.", len(numbers.Numbers))
		}
	case "stdin":
		{
			numbers, err = readFromStdIn()
			if err != nil {
				errorLog.Fatalf("Шаг 1: Чтение массива чисел из командной строки. %v", err)
			}
			infoLog.Printf("Шаг 1: Чтение массива чисел из командной строки. Чисел: %d.", len(numbers.Numbers))
		}
	default:
		errorLog.Fatalf("Пожалуйста, укажите либо 'file', либо 'stdin'")
	}

	// Задание 2.2  Сумма всех чисел в массиве.
	sum := sumArray(numbers.Numbers)
	infoLog.Printf("Шаг 2: Суммирование всех чисел в массиве. Сумма: %d.", sum)

	//Задание 2.3 Выполняет HTTP GET запрос на заданный URL и проверяет статус ответа (должен быть 200).
	if err := godotenv.Load("app.env"); err != nil {
		errorLog.Fatalf("Отсутствует env файл.")
	}

	url, exists := os.LookupEnv("URL")
	if !exists {
		errorLog.Fatalf("Отсутствует переменная окружения url.")
	}

	status, err := checkStatus(url)
	if err != nil {
		errorLog.Fatalf("Шаг 3: Выполнение HTTP GET запроса. %v", err)
	}
	infoLog.Printf("Шаг 3: Выполнение HTTP GET запроса. Статус: %d.", status)

	infoLog.Printf("Программа завершена успешно.")
}

func sumArray(nmb []int) int {
	sum := 0
	for _, num := range nmb {
		sum += num
	}
	return sum
}

func readJSONFromFile(filePath string) (*NumberArray, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении файла: %v.", err)
	}

	var numbers NumberArray
	err = json.Unmarshal(data, &numbers)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении файла: %v.", err)
	}

	return &numbers, nil
}

func readFromStdIn() (*NumberArray, error) {
	var numbers []int
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Ошибка при чтении ввода: %v", err)
		}

		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err != nil {
			break
		}

		numbers = append(numbers, num)
	}

	result := NumberArray{Numbers: numbers}
	return &result, nil
}

func checkStatus(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при выполнении Get запроса: %v.", err)
	}

	defer resp.Body.Close()
	status := resp.StatusCode
	if status != http.StatusOK {
		return 0, fmt.Errorf("HTTP статус не равен 200 OK. Статус: %d.", status)
	}

	return status, nil
}
