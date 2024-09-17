package main

import (
	"net/http"
	"fmt"
)

func main(){
	mux := http.NewServeMux()
	s := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil{
		fmt.Println("Ecounter Error during Server Initiation:", err)
		return
	}
}
