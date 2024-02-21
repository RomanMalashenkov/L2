package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	//наша строка для распаковки
	a := "a4abc2d5e"
	//основная функция распаковки
	fmt.Println(unpack(a))
}

func unpack(stroke string) string {

	//задаем переменные: строка корректна, массив цифр в строке, новая строка в которую будут добавляться данные
	var (
		ok, digits = getStrokeDigits(stroke)
		newStroke  = ""
	)

	//провермяем что строка корректна
	if !ok {
		log.Fatalln("некорректная строка")
	}

	//разбиваем строку на массив "слов", состоящих только из букв
	parts := regexp.MustCompile(`[0-9]`).Split(stroke, -1)

	//проходимся по массиву цифр в строке, сопоставляя индексы цифр и слов
	for i := 0; i < len(digits); i += 1 {
		partLen := len(parts[i])

		//проверяем что слово состоит из одной буквы
		if partLen > 1 {
			//если нет, то сначала в новую строку добавляем все символы кроме последнего, а затем последний символ
			//необходимое кол-во раз
			newStroke += parts[i][:partLen-1]
			newStroke += strings.Repeat(parts[i][partLen-1:], digits[i])
		} else {
			//иначе просто добавляем слово определенное число раз
			newStroke += strings.Repeat(parts[i], digits[i])
		}
	}

	//добавляем в строку "остатки" строки, которые содержатся после цифр
	for _, i := range parts[len(digits):] {
		newStroke += i
	}

	return newStroke
}

// функция проверки валидности строки и в случае валидности возвращает массив цифр, встреч. в исходной строке
func getStrokeDigits(stroke string) (bool, []int) {

	//компилим регексы для нахождения в строке букв и цифр
	reForNums, _ := regexp.Compile(`[0-9]`)
	reForChars, _ := regexp.Compile(`[A-z]`)

	var (
		//переменная со всеми цифрами, но в формате string
		allDigits = reForNums.FindAllString(stroke, -1)
		//количесвто цифр в строке
		counter = len(allDigits)
		//количество букв в строке
		counterForChars = len(reForChars.FindAllString(stroke, -1))
		//массив возвращаемых цифр с заранее выделенной памятью под
		//кол-во наших цифр (экономим ресурсы перевыделения памяти)
		retArr = make([]int, 0, counter)
	)

	//проверяем что строка состоит не только из цифр
	if counter > 0 && counterForChars == 0 {
		return false, nil
	}

	//добавляем в возвращаемый массив все цифры, переводя из формата string в формат int
	for _, i := range allDigits {
		//конвертим стринг цифру в инт с отлавливанием ошибки
		prev, err := strconv.Atoi(i)

		//паникуем если не получилось конвертнуть строчную цифру
		if err != nil {
			panic(err)
		}

		//если все ок - добавляем в возвращаемый массив нашу циферку
		retArr = append(retArr, prev)
	}

	return true, retArr
}
