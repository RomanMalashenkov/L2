package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type input struct {
	command, path string
	flags         []string
}

func main() {

	//создаем буфер прослушки вводимых в консоль даныных
	reader := bufio.NewReader(os.Stdin)

	//логаем изначальный заголовок
	fmt.Println("Sort Shell")
	fmt.Println("---------------------")

	//слушаем консоль в бесконечном цикле
	for {
		fmt.Print("command> ")

		//слушаем строку из буфера, деленной по \n
		buf, _ := reader.ReadString('\n')

		//отделяем нулевой элемент вводимой строки от всех отсальных, это будет команда,
		//в нашем случае либо sort либо exit. Также триммируем суффиксы для получения
		//конечной команды в виде строки
		command := strings.Split(buf, " ")[0]
		command = strings.TrimSuffix(command, "\n")
		command = strings.TrimSuffix(command, "\r")

		//завершаем программу если команда exit
		if command == "exit" {
			os.Exit(0)
		}

		//проверяем что введена и команда и второй аргумент, флаги опциональны
		if len(strings.Split(buf, " ")) < 2 {
			fmt.Println("Wrong arguments")
			continue
		}

		//создаем новый элемент команды
		returnable := &input{}

		//ставим команду элемента
		returnable.command = command
		//ставим путь элемента
		returnable.path = strings.Split(buf, " ")[1]

		//парсим флаги по аналогии с командой, триммируя \n и \r
		for _, i := range strings.Split(buf, " ")[2:] {
			i = strings.TrimSuffix(i, "\n")
			i = strings.TrimSuffix(i, "\r")
			returnable.flags = append(returnable.flags, i)
		}

		//триммируем путь так же
		returnable.path = strings.TrimSuffix(returnable.path, "\n")
		returnable.path = strings.TrimSuffix(returnable.path, "\r")

		//парсим наш элемент
		returnable.parseInput()
	}
}

// парсер команды
func (inp *input) parseInput() {

	//проверяем если команда не сорт - значит неизвестная
	if inp.command != "sort" {
		fmt.Println("Unknown command")
		return
	}

	//для наилучешго взаимодействия проинициализируем глобальный путь к исполнаяемому файлу
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	//склеиваем путь из консоли и глобальный путь
	mainPath := filepath.Join(path, inp.path)

	//получаем список наших слов
	words, errRead := ReadFile(mainPath)
	if errRead != nil {
		return
	}

	//открываем файл по глобальному пути с флагами о перезаписи данных
	cont, errParse := os.OpenFile(mainPath, os.O_WRONLY|os.O_CREATE, 0666)
	if errParse != nil {
		fmt.Printf("Error at opening file: %v\n", errParse)
		return
	}

	//по окончании функции закрываем файл
	defer cont.Close()

	//делаем простую сортировку если флагов нет
	if len(inp.flags) == 0 {

		//создаем 2 переменные: для временного копирования наших слов,
		//и двумерный массив байт для записи
		var (
			justWrite       = make([]string, len(words))
			finalJustString = make([][]byte, 0, len(words))
		)

		//копируем данные в объявленный массив и сортируем
		copy(justWrite, words)
		sort.Strings(justWrite)

		//переводим каждое слово в массив байт
		//и добавляем в конечный массив для записи
		for _, i := range justWrite {
			finalJustString = append(finalJustString, []byte(i+"\n"))
		}

		//добавляем в файл наши слова
		writeToFile(finalJustString, cont)

		//если все ок, выводим сообщение и выходим из функции
		fmt.Println("file sorted")
		return
	}

	//если же флаги есть, поочередно парсим каждый
	switch inp.flags[0] {

	//парсим флаг для колонок
	case "-k":
		//проверяем что задано значение для флага
		if inp.flags[1] == "0" {
			fmt.Println("Incorrect flag argument")
			return
		}

		//конвертим string индекс в int
		index, errConv := strconv.Atoi(inp.flags[1])
		if errConv != nil {
			fmt.Println("Cannot convert char to int")
			return
		}

		//создаем массив для сортировки строк по колонке
		//и двумерный массив из слов для каждой строки
		var (
			sortedArray    = make([]string, 0, len(words))
			finalCopyArray = make([][]string, len(words))
		)

		//для экономии ресурсов копируем в конечный массив
		//конвертнутые данные из исходного массива строк файла
		copy(finalCopyArray, convertToDoubleArray(words))

		//формируем массив слов по индексу для каждой строки файла
		for _, i := range words {
			tempWord := getWordById(index, i)
			sortedArray = append(sortedArray, tempWord)
		}

		//сортируем полученные данные
		sort.Strings(sortedArray)

		//заменяем в конечном массиве уже отсортированные слова по необходимому индексу
		for i := 0; i < len(finalCopyArray); i += 1 {
			finalCopyArray[i][index] = sortedArray[i]
		}

		//перезаписываем готовый массив в файл
		writeToFile(doubleStringToDoubleByte(finalCopyArray), cont)
	case "-r":
		//создаем 2 переменные: для временного копирования наших слов,
		//и двумерный массив байт для записи
		var (
			justWrite       = make([]string, len(words))
			finalJustString = make([][]byte, 0, len(words))
		)

		//копируем данные в объявленный массив и сортируем
		copy(justWrite, words)
		sort.Strings(justWrite)

		//делаем реверсивный массив
		justWrite = reverseStringArray(justWrite)

		//переводим каждое слово в массив байт
		//и добавляем в конечный массив для записи
		for _, i := range justWrite {
			finalJustString = append(finalJustString, []byte(i+"\n"))
		}

		//добавляем в файл наши слова
		writeToFile(finalJustString, cont)

		//если все ок, выводим сообщение и выходим из функции
		fmt.Println("file sorted")
		return
	}
	return
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

// получаем слово из по id
func getWordById(id int, words string) string {
	wds := strings.Split(words, " ")
	if len(wds)-1 < id {
		return wds[len(wds)-1]
	}
	return wds[id]
}

// разбиваем массив строк в двумерный массив из слов исходного массива
func convertToDoubleArray(words []string) [][]string {
	tempArr := make([][]string, 0, len(words))

	for _, i := range words {
		tempArr = append(tempArr, strings.Split(i, " "))
	}

	return tempArr
}

// конвертим двумерный массив строк в двумерный массив байт
func doubleStringToDoubleByte(words [][]string) [][]byte {
	returnable := make([][]byte, 0, len(words))
	for _, i := range words {
		tempStroke := ""
		for _, j := range i {
			tempStroke += j + " "
		}
		strings.TrimSuffix(tempStroke, " ")
		tempStroke += "\n"
		returnable = append(returnable, []byte(tempStroke))
	}
	return returnable
}

// пишем данные в файл
func writeToFile(finalJustString [][]byte, cont *os.File) {
	//пословно пишем каждый массив байт в строку файла
	for _, i := range finalJustString {
		_, erWrite := cont.Write(i)
		//парсим ошибку добавления, если она есть
		if erWrite != nil {
			fmt.Printf("Error at wtiting string %s to file: %v\n", string(i), erWrite)
		}
	}
}

// создаем массив в обратном порядке
func reverseStringArray(arr []string) []string {
	returnable := make([]string, 0, len(arr))
	for i := len(arr) - 1; i >= 0; i -= 1 {
		returnable = append(returnable, arr[i])
	}
	return returnable
}

//тестовые варианты для файла

//jh
//vb
//nv
//bn
//vd
//as
//fd
//gdg
//dfg
//f

//gdg gdg
//jh jh
//nv nv
//vb vb
//vd vd
//as as
//bn bn
//dfg dfg
//f f
//fd fd
