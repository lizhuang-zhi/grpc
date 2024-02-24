package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var usageString = fmt.Sprintf(`Usage: %s <integer> [-h|-help]

A greeter application which prints the name you entered <integer> number of times.
`, os.Args[0])

func printUsage(w io.Writer) {
	fmt.Fprintf(w, usageString)
}

// 获取输入的名字
func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprintf(w, msg) // 输出提示信息

	scanner := bufio.NewScanner(r)
	scanner.Scan()                        // Scan函数的默认行为是在读取换行符后返回
	if err := scanner.Err(); err != nil { // Err方法用于获取扫描过程中的错误
		return "", err
	}

	name := scanner.Text() // 获取输入的名字
	if len(name) == 0 {
		return "", errors.New("name cannot be empty")
	}

	return name, nil
}

type config struct {
	numTimes   int  // 执行次数
	printUsage bool // 是否打印用法
}

// 解析参数
func parseArgs(args []string) (config, error) {
	var numTimes int
	var err error
	c := config{}

	if len(args) != 1 {
		return c, errors.New("invalid number of arguments")
	}

	if args[0] == "-h" || args[0] == "--help" {
		c.printUsage = true
		return c, nil
	}

	numTimes, err = strconv.Atoi(args[0])
	if err != nil {
		return c, err
	}

	c.numTimes = numTimes
	return c, nil
}

// 验证参数
func validateArgs(config config) error {
	if config.printUsage {
		return errors.New("print usage")
	}
	if config.numTimes < 1 {
		return errors.New("number of times must be a positive integer")
	}
	return nil
}

// 打印名字
func greetUser(w io.Writer, name string, num int) {
	msg := fmt.Sprintf("Hello, %s!\n", name)
	for i := 0; i < num; i++ {
		fmt.Fprintln(w, msg)
	}
}

// 执行
func runCmd(r io.Reader, w io.Writer, c config) error {
	if c.printUsage {
		printUsage(w) // 打印用法
		return nil
	}

	// 获取输入的名字
	name, err := getName(r, w)
	if err != nil {
		return err
	}

	// 打印名字
	greetUser(w, name, c.numTimes)
	return nil
}

func main() {
	c, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		printUsage(os.Stdout)
		os.Exit(1) // 退出并返回错误码
	}

	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		printUsage(os.Stdout)
		os.Exit(1)
	}

	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
