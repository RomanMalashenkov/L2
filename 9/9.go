package main

//эта программа выполняет скачивание файла по адресу.
//Скачивание сайта целиком описано в файле wget_website

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type object struct {
	filename, href string
	wFileName      bool
	params         [][]string
	limit          int
}

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

		//получаем ссылку из введенного текста
		href := splited[len(splited)-1]
		href = strings.TrimSuffix(href, "\n")
		href = strings.TrimSuffix(href, "\r")

		//задаем параметры если они есть
		params := make([][]string, 0)
		var prm []string
		if len(splited) > 2 {
			prm = splited[1 : len(splited)-2]
		}

		//закидываем параметры по парам флаг - значение
		flagMustSkip := false
		for i := 0; i < len(prm); i += 1 {
			if flagMustSkip {
				flagMustSkip = false
			} else {
				params = append(params, prm[i:i+2])
				flagMustSkip = true
			}
		}

		//создем объект скачивания и качаем его
		ob := object{
			href:   href,
			params: params,
		}
		err := ob.download()
		if err != nil {
			fmt.Println(err)
		}

	}
}

// скачиваем объект по его данным
func (o *object) download() error {
	//получаем файл по ссылке
	file, errGet := http.Get(o.href)
	if errGet != nil {
		fmt.Println("Error at getting file: ", errGet)
		return errGet
	}

	//парсим флаги в параметрах объекта
	for _, flag := range o.params {
		if flag[0] == "-O" {
			o.wFileName = true
			o.filename = flag[1]
		}
		if strings.Split(flag[0], "=")[0] == "––limit-rate" {
			o.limit, _ = strconv.Atoi(strings.Split(flag[0], "=")[1])
		}
	}

	defer file.Body.Close()

	//создаем конечную переменную скачанных данных,
	//буфер обмена,
	//последнтие 2 параметра необходимы для расчета скоростит скачивания
	var downloadedBytes int
	buf := make([]byte, 1024)
	prevTime := time.Now()
	prevDownloadedBytes := 0

	//итерируемся по целевому файлу, пока не скачаем его полностью
	for {
		//читаем буфер целевого файли и пшем в наш
		n, err := file.Body.Read(buf)
		//отсчитали прошедшее время
		currentTime := time.Now()
		//считаем скорость скачивания, процент скачки
		elapsedTime := currentTime.Sub(prevTime).Seconds()
		downloadSpeed := float64(downloadedBytes-prevDownloadedBytes) / 1024 / elapsedTime
		currentString := [100]byte{}
		percentsFloat := float64(downloadedBytes) / float64(file.ContentLength) * 100
		percents := int(percentsFloat)
		//заполняем строку процессинга
		for i := 0; i < percents; i += 1 {
			currentString[i] = '#'
		}
		for i := percents; i < 100; i += 1 {
			currentString[i] = '.'
		}
		//если буфер для скачки не пустой - продалжаем качать
		if n > 0 {
			downloadedBytes += n
			fmt.Printf("\r%d of %d KB downloaded |%s| %d %s | Speed : %.2f KB/s", downloadedBytes/1024, file.ContentLength/1024, currentString, percents, "%", downloadSpeed)
		}
		if err != nil {
			break
		}
	}

	//строка состояния скачания
	fmt.Println("Download successful\nFinal file size:", downloadedBytes/1024, "KB")
	defer file.Body.Close()

	//если имя файла не установлено - ставим название сайта
	if !o.wFileName {
		o.filename = path.Base(file.Request.Host)
	}

	//создаем файл с готовым именем
	cr, errCreate := os.Create(o.filename)
	if errCreate != nil {
		fmt.Println("Error at creating file: ", errCreate)
		return errCreate
	}

	defer cr.Close()
	//копируем данные в файл
	_, errWrite := io.Copy(cr, file.Body)
	if errWrite != nil {
		fmt.Println("Error at writing at file: ", errWrite)
		return errWrite
	}
	return nil
}
