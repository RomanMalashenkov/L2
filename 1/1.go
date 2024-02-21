package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {

	//получаем время с сервера ntp
	timeNtp, err := ntp.Time("pool.ntp.org")

	//парсим ошибки
	if err != nil {
		b, _ := fmt.Fprintf(os.Stderr, "Error at getting time: %v\n", err)
		fmt.Printf("%d bytes written", b)
		os.Exit(1)
	}

	//выводим время
	fmt.Printf("Current time: %s", timeNtp.Format(time.RFC3339))

}
