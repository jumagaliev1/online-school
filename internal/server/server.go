package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ServeHTTP(r *mux.Router, port string) {
	fmt.Println("starting in: ", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println(err)
	}

}
