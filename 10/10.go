package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type connection struct {
	address, port string
	socket        net.Conn
	params        []string
}

func main() {
	//создаем буфер прослушки вводимых в консоль даныных
	reader := bufio.NewReader(os.Stdin)

	//логаем изначальный заголовок
	fmt.Println("Telnet Shell")
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
		if len(splited) < 3 {
			fmt.Println("Wrong arguments")
			continue
		}

		//задаем переменные адреса и порта исходя из параметров
		address := splited[1]
		port := splited[2]

		//тримируем порт тк он последний в строке
		port = strings.TrimSuffix(port, "\n")
		port = strings.TrimSuffix(port, "\r")

		//парсим флаги по аналогии с командой, триммируя \n и \r
		var tempFlags []string
		for _, i := range strings.Split(buf, " ")[3:] {
			i = strings.TrimSuffix(i, "\n")
			i = strings.TrimSuffix(i, "\r")
			tempFlags = append(tempFlags, i)
		}

		//создаем объект нашего соединения и запускаем стартовую функцию
		conn := &connection{address: address, port: port, params: tempFlags}
		conn.start()
	}
}

// запускаем все процессы и проверки
func (c connection) start() {
	//создаем адрес
	address := fmt.Sprintf("%s:%s", c.address, c.port)
	fmt.Println("Trying", address, "...")
	//таймаут на подключение, по дефолту - 0
	//если флаг таймаута задан - парсим его
	timeOut := 0 * time.Second
	if len(c.params) != 0 {
		prvTime := strings.TrimSuffix(strings.Split(c.params[0], "=")[1], "s")
		prvTimeInt, _ := strconv.Atoi(prvTime)
		timeOut = time.Duration(prvTimeInt) * time.Second
	}

	//создаем соединение
	conn, errConnect := net.DialTimeout("tcp", address, timeOut)
	if errConnect != nil {
		fmt.Println("Err at read stdin:", errConnect)
		os.Exit(1)
	}

	//ставим переменную сокета в объекте подключения
	c.socket = conn

	fmt.Println("Connected to", address)

	//параллельно запускаем чтение из сокета и чтение с консоли
	go read(c)
	//go listen(c, &wg)

	buf := make([]byte, 8192)
	for {
		//читаем Stdin, кидаем в буфер
		fmt.Print("telnet> ")
		inp, errRead := os.Stdin.Read(buf)
		if errRead != nil {
			fmt.Println("Err at read stdin:", errRead)
			os.Exit(1)
		}

		//пишем полученные данные в сокет
		_, errSockWrite := c.socket.Write(buf[:inp])
		if errSockWrite != nil {
			fmt.Println("Err at write socket:", errSockWrite)
			os.Exit(1)
		}
	}
}

// чтение из сокета
func read(c connection) {
	//буфер получаемых данных
	buf := make([]byte, 8192)
	for {
		//слушаем сокет и выводим полученные данные
		inp, err := c.socket.Read(buf)
		if err != nil {
			fmt.Println("Err at read socket:", err)
			os.Exit(1)
		}
		fmt.Println(string(buf[:inp]))
	}
}

// чтение с консоли
func listen(c connection) {
	//буфер вводимых данных
	buf := make([]byte, 1024)
	for {
		//читаем Stdin, кидаем в буфер
		fmt.Print("telnet> ")
		inp, errRead := os.Stdin.Read(buf)
		if errRead != nil {
			fmt.Println("Err at read stdin:", errRead)
			os.Exit(1)
		}

		//пишем полученные данные в сокет
		_, errSockWrite := c.socket.Write(buf[:inp])
		if errSockWrite != nil {
			fmt.Println("Err at write socket:", errSockWrite)
			os.Exit(1)
		}
	}
}
