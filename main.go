package main

import (
    "main/db"
    "main/server"
    "fmt"
)

func main() {
    fmt.Println("hello world")
    //go server.StartServer()
    db.Start()
    server.StartHttp()
}
