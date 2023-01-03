package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var path string

func init() {
	flag.StringVar(&path, "p", "", "please input path")
}

func main() {
	//CheckDir("test/")
	CheckDir("gamesvr/")
}

func Replace(file string) error {
	fmt.Println("file: ", file)

	fr, err := os.Open(file)
	if err != nil {
		fmt.Printf("os.Open(%v) failed(%v)", file, err)
		return nil
	}
	buf := bufio.NewReader(fr)

	lines := []string{}
	for {
		//line, err := buf.ReadString('\n')
		byteline, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		lines = append(lines, string(byteline)+"\n")
	}
	fr.Close()
	if len(lines) <= 0 {
		return fmt.Errorf("%v file is empty", file)
	}

	fnew, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	//f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		fmt.Printf("os.Open(%v) failed(%v)", file, err)
		return nil
	}
	defer fnew.Close()

	for _, v := range lines {
		newLine := v

		if strings.Contains(v, "Tables()") &&
			!strings.Contains(v, "ecdoe.Errorf") &&
			!strings.Contains(v, "ecode.ErrMsgf") &&
			!strings.Contains(v, "logger.") {
			newLine = ""

			index := strings.Index(v, "Tables().")
			if index > 0 {
				right := v[0:index]
				left := v[(index + len("Tables().")):]
				leftIndex := strings.Index(left, ".")
				newLine = right + left[0:leftIndex] + "()" + left[leftIndex:]
				//fmt.Println("czx@@@ newLine: ", newLine)
			}
		}
		fnew.WriteString(newLine)
	}
	fnew.Close()

	return nil
}

func CheckDir(pathname string) error {
	fmt.Println("pathname : ", pathname)
	rd, outErr := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			err := CheckDir(pathname + fi.Name() + "/")
			if err != nil {
				outErr = err
			}

		} else {
			err := Replace(pathname + fi.Name())
			if err != nil {
				outErr = err
			}
		}
	}

	return outErr
}
