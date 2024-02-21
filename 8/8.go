package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

type proc struct {
	route string
}

func main() {
	//создаем буфер прослушки вводимых в консоль даныных
	reader := bufio.NewReader(os.Stdin)

	//логаем изначальный заголовок
	fmt.Println("Unix Shell")
	fmt.Println("---------------------")

	dir, errCurPath := os.Getwd()

	if errCurPath != nil {
		fmt.Printf("Cannot resolve path: %v\n", errCurPath)
	}

	main := &proc{route: dir}

	//слушаем консоль в бесконечном цикле
	for {
		fmt.Printf("Shell | %s> ", main.route)

		//слушаем строку из буфера, деленной по \n
		buf, _ := reader.ReadString('\n')

		//отделяем нулевой элемент вводимой строки от всех отсальных, это будет команда,
		//в нашем случае либо sort либо exit. Также триммируем суффиксы для получения
		//конечной команды в виде строки
		command := strings.Split(buf, " ")[0]
		command = strings.TrimSuffix(command, "\n")
		command = strings.TrimSuffix(command, "\r")

		//завершаем программу если команда \quit
		if command == "\\quit" {
			os.Exit(0)
		}

		//проверяем что введена команда
		if len(strings.Split(buf, " ")) < 1 {
			fmt.Println("Wrong arguments")
			continue
		} else if len(strings.Split(buf, " ")) == 1 {
			//если введена только команда, без аргументов - ставим массив аргументов нулевым
			main.parseCommand(command, nil)
		} else {
			//иначе парсим аргументы и триммируем последний
			args := strings.Split(buf, " ")[1:]
			args[len(args)-1] = strings.TrimSuffix(args[len(args)-1], "\n")
			args[len(args)-1] = strings.TrimSuffix(args[len(args)-1], "\r")
			main.parseCommand(command, args)
		}
	}
}

// основная функция парсинга флагов
func (p *proc) parseCommand(command string, args []string) {
	switch command {
	case "echo":
		//если команда - эхо - просто выводим все аргументы
		printClearArray(args)
	case "cd":
		//работаем с окружением
		//если не введено аргументов - не выполняем никаких действий
		if len(args) != 0 {
			//парсим "выход вверх" текущей директории
			if args[0] == ".." {
				p.route = routeUp(p.route)
			} else {
				//иначе преобразуем путь в массив рун для корректной работы
				w := []rune(args[0])
				//если путь глобальный - проверяем его существование и ставим текущую директорию туда
				if string(w[0]) == "/" {
					file, err := absPathErr(args[0])
					if err == nil {
						p.route = file
					} else {
						fmt.Println(err)
					}
				} else {
					//если же путь локальный - проверяем существование внутри
					//текущей директории и если все ок - переходим туда
					file, err := tryGetPath(p.route, args[0])
					if err == nil {
						p.route = file
					} else {
						fmt.Println(err)
					}
				}
			}
		}
	case "pwd":
		//логаем текущую директорию
		p.logCurrentRoute()
	case "ps":
		//логаем все запущенные процессы
		logProcess()
	case "kill":
		//убиваем процесс по PID, если такой есть
		Pid, errConv := strconv.Atoi(args[0])
		if errConv != nil {
			fmt.Println("Err at converting:", errConv)
			return
		}
		errKill := killProcById(Pid)
		if errKill != nil {
			fmt.Println("Err at killing:", errKill)
			return
		}
		fmt.Println("Process", args[0], "killed")
	case "ls":
		//логаем окружение в текущей директории
		p.getAndLogEnv()
	}

}

// "чистый" лог - массив выводится в одну строку как обычная строка
func printClearArray(args []string) {
	for _, i := range args {
		fmt.Print(i + " ")
	}
	fmt.Println()
}

// поднимаемся на папку вверх
func routeUp(currentRoute string) string {
	//разбиваем роут по \  и получаем список папок,
	//создавая новый путь без последнего элемента
	splitedRoute := strings.Split(currentRoute, "\\")
	newRoute := splitedRoute[:len(splitedRoute)-1]

	//билдим наш путь в едину строку из массива
	//если текущая папка является "последней" - дальше не идем
	builder := concString(newRoute)
	if len(strings.Split(builder, "\\")) <= 2 {
		return builder
	}
	return strings.TrimSuffix(builder, "\\")
}

// проверяем и получаем абсолютный путь
func absPathErr(path string) (string, error) {
	if _, err := os.Stat(path); err == nil {
		pInfo, err := filepath.Abs(path)
		return pInfo, err
	} else {
		return "", err
	}
}

// складываем пути, проверяем что они вообще есть
// (функция для локального взаимодействия)
func tryGetPath(oldPath, path string) (string, error) {
	newPath := filepath.Join(oldPath, path)
	if _, err := os.Stat(newPath); err == nil {
		return newPath, nil
	} else {
		return "", err
	}
}

// простая функция конкатенации строк
func concString(arr []string) string {
	var builder strings.Builder
	for _, s := range arr {
		_, err := builder.WriteString(s + "\\")
		if err != nil {
			log.Fatal(err)
		}
	}
	return builder.String()
}

// получаем и логаем окружение
func (p *proc) getAndLogEnv() {
	all, err := os.ReadDir(p.route)
	if err != nil {
		fmt.Println("Error at parsing enviroment:", err)
		return
	}
	for _, obj := range all {
		fmt.Println(obj)
	}
}

// получаем и логаем запущенные процессы с PID и названием
func logProcess() {
	processes, err := ps.Processes()
	if err != nil {
		fmt.Println("Err at logging processes:", err)
		return
	}
	fmt.Println("   PID\t| Executable")
	for i := range processes {
		proc := processes[i]
		fmt.Printf("%d\t| %s\n", proc.Pid(), proc.Executable())
	}
}

// убиваем процесс по PID встроенными методами языка
func killProcById(Pid int) error {
	proc, errFind := os.FindProcess(Pid)
	if errFind != nil {
		return errFind
	}
	errKill := proc.Kill()
	if errFind != nil {
		return errKill
	}
	return nil
}

// логаем текущий путь в формате для команды pwd
func (p *proc) logCurrentRoute() {
	fmt.Println()
	fmt.Println("Path")
	fmt.Println("----")
	fmt.Println(p.route)
	fmt.Println()
}
