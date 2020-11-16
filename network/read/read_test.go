package connread

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"testing"
)

func BenchmarkRead_Once(b *testing.B) {
	b.StopTimer()

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bufReader := bufio.NewReaderSize(f, 1024)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := Read_Bufio(bufReader)
		if err != nil {
			panic(err)
		}
		_, err = Read_Bufio(bufReader)
		if err != nil {
			panic(err)
		}
		bufReader.Reset(f)
		syscall.Seek(int(f.Fd()), 0, 0)
	}
}

func BenchmarkRead_Once2(b *testing.B) {
	b.StopTimer()

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := make([]byte, 1024)
	bufReader := &ByteBuf{
		buf:  bytes.NewBuffer(nil),
		size: -1,
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n, err := f.Read(buf)
		if err != nil {
			log.Println(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}
		_, err = bufReader.Read_ByteBuf(buf[:n])
		if err != nil {
			panic(err)
		}
		syscall.Seek(int(f.Fd()), 0, 0)
	}
}

func BenchmarkRead_Twice(b *testing.B) {
	b.StopTimer()

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := Read_ReadFull(f)
		if err != nil {
			panic(err)
		}
		_, err = Read_ReadFull(f)
		if err != nil {
			panic(err)
		}
		syscall.Seek(int(f.Fd()), 0, 0)
	}
}

func _TestRead(t *testing.T) {
	{
		f, err := os.Open("data")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		bufReader := bufio.NewReaderSize(f, 1024)
		data, err := Read_Bufio(bufReader)
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(data), "message:", string(data))

		data, err = Read_Bufio(bufReader)
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(data), "message:", string(data))
	}

	{

		f, err := os.Open("data")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buf := make([]byte, 1024)
		bufReader := &ByteBuf{
			buf:  bytes.NewBuffer(nil),
			size: -1,
		}

		n, err := f.Read(buf)
		if err != nil {
			log.Println(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}
		bufs, err := bufReader.Read_ByteBuf(buf[:n])
		if err != nil {
			panic(err)
		}
		for _, v := range bufs {
			log.Println("length:", len(v), "message:", string(v))
		}
	}

	// copy once, call sys.read twice
	{
		f, err := os.Open("data")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		content, err := Read_ReadFull(f)
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(content), "message:", string(content))

		content, err = Read_ReadFull(f)
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(content), "message:", string(content))
	}
}
