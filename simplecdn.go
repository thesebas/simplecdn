package main

import (
	"fmt"
	redis "github.com/tidwall/redcon"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	done := make(chan bool)
	go func() {
		log.Println("listening for http connections")

		err := http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(httpServe))
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Println("listening for redis connections")
		err := redis.ListenAndServe("0.0.0.0:6379", redisServe, redisAccept, redisClosed)
		if err != nil {
			log.Fatal(err)
		}

	}()
	<-done
}

func httpServe(w http.ResponseWriter, r *http.Request) {
	log.Printf("http requested %s %s", r.Method, r.RequestURI)
	w.WriteHeader(404)
}

func redisServe(conn redis.Conn, cmd redis.Command) {
	args := cmd.Args
	log.Println(args)
	command := args[0]

	for idx, str := range args {
		fmt.Println(idx, string(str))
	}

	switch strings.ToLower(string(command)) {
	default:

		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "mget":
		files := args[1:]

		conn.WriteArray(len(args) - 1)

		for idx, file := range files {
			conn.WriteString("http://127.0.0.1:8080/somefile/" + string(file) + "?idx=" + strconv.Itoa(idx))
		}

	}
}

func redisAccept(conn redis.Conn) bool {
	log.Printf("redis accepted %s", conn.RemoteAddr())
	return true
}

func redisClosed(conn redis.Conn, err error) {
	if err != nil {
		log.Printf("redis unexpected closed %s %s", conn.RemoteAddr(), err.Error())
	} else {
		log.Printf("redis connection closed %s", conn.RemoteAddr())
	}
}
