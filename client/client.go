package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	conn, e := net.Dial("tcp", ":1833")
	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	go func() { //waiting response from server
		buf := bufio.NewReader(conn)
		for { //do something with response
			time.Sleep(100 * time.Millisecond)
			msg, err := buf.ReadBytes('\n')
			if err != nil {
				panic(err)
			}

			fmt.Printf("Message from server: %s", string(msg))
		}
	}()

	f, err := os.Open("book.txt") // data for example
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	buf := make([]byte, 0, 2000)
	for {
		time.Sleep(2 * time.Second)
		randLen := rand.Intn(500) // data will be sent in chunks of 1-500 bytes
		fmt.Printf("take %d bytes from file\n", randLen)
		n, err := io.ReadFull(r, buf[:randLen])
		buf = buf[:n]
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.ErrUnexpectedEOF {
				log.Println(err)
				break
			}
		}
		_, err = conn.Write(buf) // send data to server
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	// -=Old=-
	//example := [][]byte{
	//	[]byte("111\n222\n333"),
	//	[]byte("444"),
	//	[]byte("777\n"),
	//	[]byte("\n"),
	//	[]byte("\n666\n777"),
	//	[]byte("111\n222"),
	//	[]byte("111\n222\n"),
	//}
	//
	//for _, bytes := range example {
	//	time.Sleep(2 * time.Second)
	//	fmt.Println("Sent to server =>", string(bytes))
	//	_, err := conn.Write(bytes)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	select {}

}

//package main
//
//import (
//	"bufio"
//	"encoding/binary"
//	"fmt"
//	"io"
//	"log"
//	"net"
//	"os"
//)
//
//func Packet(data []byte) []byte {
//	buffer := make([]byte, 4+len(data))
//	// 将buffer前面四个字节设置为包长度，大端序
//	binary.BigEndian.PutUint32(buffer, uint32(len(data)))
//	copy(buffer[4:], data)
//	return buffer
//}
//
//func UnPacket(c net.Conn) ([]byte, error) {
//	var header = make([]byte, 4)
//
//	_, err := io.ReadFull(c, header)
//	if err != nil {
//		return nil, err
//	}
//	length := binary.BigEndian.Uint32(header)
//	contentByte := make([]byte, length)
//	_, e := io.ReadFull(c, contentByte) //读取内容
//	if e != nil {
//		return nil, e
//	}
//
//	return contentByte, nil
//}
//
//func main() {
//	conn, e := net.Dial("tcp", ":1833")
//	if e != nil {
//		log.Fatal(e)
//	}
//	defer conn.Close()
//
//	for {
//		reader := bufio.NewReader(os.Stdin)
//		fmt.Print("Text to send: ")
//		text, _ := reader.ReadString('\n')
//
//		buffer := Packet([]byte(text))
//		_, err := conn.Write(buffer)
//		if err != nil {
//			panic(err)
//		}
//
//		// listen for reply
//		msg, err := UnPacket(conn)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Printf("Message from server (len %d) : %s", len(msg), string(msg))
//	}
//}
