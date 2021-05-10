package im

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
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
		if tmpl, err := template.ParseFiles("index.html"); err != nil {
			log.Errorf("error parse file:%v", err)
			return
		} else {
			//http.ServeFile(w, r, "index.html")
			tmpl.Execute(w, map[string]interface{}{"uid": uid})
			return
		}
	}

}

func Run(httpPort int, grpcPort int) {
	flag.Parse()
	hub := newHub()
	go hub.run()
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("url:%s path:%s", r.URL, r.URL.Path))
		if uid, err := strconv.ParseInt(r.URL.Query().Get("uid"), 10, 64); err != nil {
			w.WriteHeader(400)
			w.Write([]byte("please login first"))
			return
		} else {
			serveWs(hub, w, r, uid)
		}
	})
	addr := fmt.Sprintf("0.0.0.0:%d", httpPort)
	log.Printf("listening addr:%s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
