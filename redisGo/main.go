// Creating simple example of redis key value store
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	myLog, err := os.Create("Key-Value Log")
	if err != nil {
		fmt.Println(err)
	}

	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()
	defer myLog.Write([]byte("Done for today"))

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		myLog, err := os.Create("Key-Value Log")
		go handle(conn, myLog.Name())
	}
}

func handle(conn net.Conn, file string) {
	defer conn.Close()

	// User Instructions
	io.WriteString(conn, "\r\n KEY VALUE DB\r\n\r\n"+
		"USE:\r\n"+
		"\tSET key value \r\n"+
		"\tGET key \r\n"+
		"\tDEL key \r\n\r\n"+
		"EXAMPLE:\r\n"+
		"\tSET fav chocolate \r\n"+
		"\tGET fav \r\n\r\n\r\n")

	data := make(map[string]string)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()     
		fs := strings.Fields(ln) 

		if len(fs) < 1 { 
			continue
		}

		switch fs[0] {
		case "GET":
			k := fs[1]                                   
			v := data[k]                                 
			fmt.Fprintf(conn, "You requested %s\r\n", v) 
		case "SET":
			if len(fs) != 3 {
				fmt.Fprintln(conn, "Expected value \r\n")
				continue
			}
			k := fs[1]  
			v := fs[2] 
			data[k] = v 
			addition := "Added value of " + data[k] + "to key" + fs[1] + "\n"
			ioutil.WriteFile(file, []byte(addition), 0666)
		case "DEL":
			if len(fs) != 2 {
				fmt.Fprintln(conn, "Expected item to delete \r\n")
				continue
			}
			k := fs[1] 
			delete(data, k)
			deletion := "Deleted " + k + "\n"
			ioutil.WriteFile(file, []byte(deletion), 0666)
		default: 
			fmt.Fprintln(conn, "Invalid command ", fs[0]+"\r\n")
			continue
		}
	}
}
