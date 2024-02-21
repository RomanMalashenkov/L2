package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type input struct {
	command, path, word string
	params              [][]string
}

func main() {

	//создаем буфер прослушки вводимых в консоль даныных
	reader := bufio.NewReader(os.Stdin)

	//логаем изначальный заголовок
	fmt.Println("Grep Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("command> ")

		//слушаем строку из буфера, деленной по \n
		buf, _ := reader.ReadString('\n')

		splited := strings.Split(buf, " ")

		command := splited[0]
		command = strings.TrimSuffix(command, "\n")
		command = strings.TrimSuffix(command, "\r")

		//завершаем программу если команда exit
		if command == "exit" {
			os.Exit(0)
		}

		//проверяем что введена и команда и второй аргумент, флаги опциональны
		if len(splited) < 3 {
			fmt.Println("Wrong arguments")
			continue
		}

		inp := &input{}

		inp.command = command

		inp.path = splited[len(splited)-1]
		inp.path = strings.TrimSuffix(inp.path, "\n")
		inp.path = strings.TrimSuffix(inp.path, "\r")

		inp.params = make([][]string, 0)
		inp.word = splited[len(splited)-2]

		prm := splited[1 : len(splited)-2]

		flagMustSkip := false

		for i := 0; i < len(prm); i += 1 {
			if prm[i] != "-i" && prm[i] != "-v" && prm[i] != "-n" {
				if flagMustSkip {
					flagMustSkip = false
				} else {
					inp.params = append(inp.params, prm[i:i+2])
					flagMustSkip = true
				}
			} else {
				inp.params = append(inp.params, prm[i:i+1])
			}
		}

		inp.searchForWord()

	}

}

func (I *input) searchForWord() {

	//проверяем если команда не grep - значит неизвестная
	if I.command != "grep" {
		fmt.Println("Unknown command")
		return
	}

	//для наилучешго взаимодействия проинициализируем глобальный путь к исполнаяемому файлу
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	//склеиваем путь из консоли и глобальный путь
	mainPath := filepath.Join(path, I.path)

	//получаем список наших слов
	words, errRead := ReadFile(mainPath)
	if errRead != nil {
		return
	}

	//ставии флаги false по дефолту
	bothRegs := false
	invert := false

	//парсим массив наших флагов
	for _, flag := range I.params {
		if flag[0] == "-i" {
			bothRegs = true
		}
		if flag[0] == "-v" {
			invert = true
		}
	}

	//генерим необходимый regex исходя из флагов и компилим его
	searchRegex := generateRegex(I.word, bothRegs, invert)
	regex, _ := regexp.Compile(searchRegex)

	//двуммерный массив для добавления строк с совпадениями
	searched := [][]string{}

	//ищем совпадения с помощью созданного regex
	for _, word := range words {
		temp := regex.FindAllString(word, -1)
		searched = append(searched, temp)
	}

	//для вывода результата создаем сборщик строк, куда будем помещать конечный результат
	var builder strings.Builder

	//циклом добавляем все совпадения в нашу строку
	for _, wds := range searched {
		for _, wd := range wds {
			_, er := builder.WriteString(strings.TrimPrefix(wd+"\n", " "))
			if er != nil {
				fmt.Println(er)
			}
		}
	}

	//далее создаем массив куда будем класть данные по флагам контекста (-a, -b, -C) и парсим эти флаги
	splitedReady := strings.Split(builder.String(), "\n")
	parsedContext := I.parseFlagsContext()

	//итерируемся по массиву готовых строк чтобы составить финальный результат
	for _, elem := range splitedReady {

		//ставим временный элемент равным элементу итерации
		tempElem := strings.TrimSuffix(elem, "\n\r")
		//если элемент - не пустая строка - работаем с ним
		if len(tempElem) > 0 {
			//итерируемся по найденым индексам временного элемента в общем тексте
			for _, el := range getIndexOf(words, tempElem) {
				//так же проверяем что индекс != 0
				if el != 0 {
					//выводим для каждого флага соответствующие строки
					if parsedContext["-b"] != 0 {
						for i := el - parsedContext["-b"]; i < el; i += 1 {
							fmt.Println(words[i])
						}
					} else if parsedContext["-C"] != 0 {
						for i := el - parsedContext["-C"]; i < el; i += 1 {
							fmt.Println(words[i])
						}
					}
					fmt.Println(words[el])
					if parsedContext["-a"] != 0 {
						for i := el + 1; i < el+parsedContext["-a"]+1; i += 1 {
							fmt.Println(words[i])
						}
					} else if parsedContext["-C"] != 0 {
						for i := el + 1; i < el+parsedContext["-C"]+1; i += 1 {
							fmt.Println(words[i])
						}
					}
				}
			}
		}
	}
}

// возвращает массив индексов найденых элементов в общем масисве
func getIndexOf(ar []string, elem string) []int {
	indexes := []int{}
	for i := 0; i < len(ar); i += 1 {
		if ar[i] == elem {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// создает regex исходя из параметров
func generateRegex(word string, bothRegs bool, revert bool) string {
	res := fmt.Sprintf(".*\\b(%s)\\b.*", word)
	if bothRegs {
		if revert {
			res = fmt.Sprintf("^(?:(?!\\b(%s|%s)\\b).)*$", strings.ToUpper(word), strings.ToLower(word))
		} else {
			res = fmt.Sprintf(".*\\b(%s|%s)\\b.*", strings.ToUpper(word), strings.ToLower(word))
		}
	} else {
		if revert {
			res = fmt.Sprintf("^(?:(?!\\b(%s)\\b).)*$", word)
		}
	}
	return res
}

// создает мапу для контекстных флагов
func (I *input) parseFlagsContext() map[string]int {
	ret := make(map[string]int, 3)
	ret["-a"] = 0
	ret["-b"] = 0
	ret["-C"] = 0
	for _, flag := range I.params {
		if flag[0] == "-a" || flag[0] == "-b" || flag[0] == "-C" {
			ret[flag[0]], _ = strconv.Atoi(flag[1])
		}
	}
	return ret

}

// чтение файла и возврат данных
func ReadFile(mainPath string) ([]string, error) {

	//открываем файл по указанному пути
	cont, errParse := os.Open(mainPath)
	if errParse != nil {
		fmt.Printf("Error at opening file: %v\n", errParse)
		return nil, errParse
	}

	//читаем данные
	scan := bufio.NewScanner(cont)

	//заранее создаем массив слов взятых из файла
	var words []string

	//добавляем содержимое файла в массив слов построчно
	for scan.Scan() {
		words = append(words, scan.Text())
	}

	//если все ок возвращаем массив слов и нулевую ошибку
	return words, nil
}
