package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/felts94/local-message-board/location"
)

type region map[string]city

type city struct {
	Posts []post `json:"posts"`
	Info  string `json:"name"`
}

type post struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Link    string `json:"link"`
}

// Echo headers and info
type Echo struct {
	Headers http.Header   `json:"headers"`
	IP      string        `json:"ip"`
	Info    location.Info `json:"info"`
}

var data map[string]region

func main() {
	setLogging()
	fmt.Printf("Starting server,  log file: info.log")

	var PORT string
	var HOST string

	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8081"
	}

	if HOST = os.Getenv("HOST"); HOST == "" {
		HOST = ""
	}

	_ = importJSONDataFromFile("data.json", &data)

	http.HandleFunc("/post", postMessage)
	http.HandleFunc("/read", read)
	http.HandleFunc("/info", info)
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	log.Printf("Starting Server %s", fmt.Sprintf("%s:%s", HOST, PORT))

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", HOST, PORT), nil)
	if err != nil {
		panic(err)
	}

}

func setLogging() {
	logfile := "/dev/stdout"
	lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

	if err != nil {
		log.Fatal("OpenLogfile: os.OpenFile:", err)
	}

	log.SetOutput(lf)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//Helper function to import json from file to map
func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}

func read(w http.ResponseWriter, r *http.Request) {
	info := location.GetUserLocation(r.RemoteAddr)
	log.Printf("%v", info)
	return
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	return
}

func info(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL)

	info := location.GetUserLocation(r.RemoteAddr)
	log.Printf("[%s] Local: %v", r.RemoteAddr, info)

	response := Echo{
		Headers: r.Header,
		IP:      r.RemoteAddr,
		Info:    info,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "close")
	w.Write(responseBytes)
}
