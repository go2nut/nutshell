package im

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("url:%s path:%s", r.URL, r.URL.Path))
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if uid, err := strconv.ParseInt(r.URL.Query().Get("uid"), 10, 64); err != nil {
		http.Error(w, "Uid not provide", http.StatusUnauthorized)
	} else {
		http.ServeFile(w, r, fmt.Sprintf("index.html?uid=%d", uid))
	}

}

func Run(httpPort int, grpcPort int) {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if uid, err := strconv.ParseInt(r.Form.Get("uid"), 10, 64); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("please login first"))
			return
		} else {
			serveWs(hub, w, r, uid)
		}
	})
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", httpPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
