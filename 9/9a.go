package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	//создаем буфер прослушки вводимых в консоль даныных
	reader := bufio.NewReader(os.Stdin)

	//логаем изначальный заголовок
	fmt.Println("Wget Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("command> ")

		//слушаем строку из буфера, деленной по \n
		buf, _ := reader.ReadString('\n')

		splited := strings.Split(buf, " ")

		command := splited[0]

		//завершаем программу если команда exit
		if command == "exit" {
			os.Exit(0)
		}

		//проверяем что введена и команда и второй аргумент, флаги опциональны
		if len(splited) < 2 {
			fmt.Println("Wrong arguments")
			continue
		}

		//ссылка для скачивания сайта будет посленим аргументом консоли
		href := splited[len(splited)-1]
		href = strings.TrimSuffix(href, "\n")
		href = strings.TrimSuffix(href, "\r")

		//массив с параметрами, если они будут
		params := make([][]string, 0)

		//заполняем параметры если они есть
		var prm []string
		if len(splited) > 2 {
			prm = splited[1 : len(splited)-2]
		}

		//итерируемся по массиву флагов, закидываем их в пары флаг - значение
		flagMustSkip := false
		for i := 0; i < len(prm); i += 1 {
			if flagMustSkip {
				flagMustSkip = false
			} else {
				params = append(params, prm[i:i+2])
				flagMustSkip = true
			}
		}

		//обрабатываем ошибку при скачивании
		err := download(href)
		if err != nil {
			fmt.Println(err)
		}

	}
}

// качаем файл по URI выкидывая ошибку, если есть
func download(uri string) error {
	//получаем файл по ссылке
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	//получаем контент с наешго запроса
	content, errRead := io.ReadAll(res.Body)
	if errRead != nil {
		return errRead
	}
	//пишем полученные данные в файл
	writeToFile(content, res.Request.Host)
	return nil
}

// пишем данные в файл
func writeToFile(finalJustString []byte, name string) {
	//создаем файл
	cont, errCreateFile := os.Create(name)
	if errCreateFile != nil {
		fmt.Printf("Error at creating file: %v\n", errCreateFile)
	}
	//пишем данные
	_, erWrite := cont.Write(finalJustString)
	if erWrite != nil {
		fmt.Printf("Error at wtiting to file: %v\n", erWrite)
	}
}
