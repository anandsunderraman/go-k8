package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "My Awesome Go App")
}

func config(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(filepath.Clean("/client-config/client_config.json"))
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("error reading config: %s", err))
	} else {
		fmt.Fprintf(w, fmt.Sprintf("config value is : %s", string(b)))
	}

}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/config", config)
}

func main() {
	fmt.Println("Go Web App Started on Port 3000")
	setupRoutes()
	http.ListenAndServe(":3000", nil)
}
