package im

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"nutshell/_example/apps/shard"
	userCli "nutshell/_example/apps/user/client"
)


func serveHome(w http.ResponseWriter, r *http.Request) {

	var homeTempl = template.Must(template.ParseFiles("index.html"))

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if token := r.URL.Query().Get("token"); token == "" {
		http.Error(w, "token not provide", http.StatusUnauthorized)
	} else {
		//if tmpl, err := template.ParseFiles("index.html"); err != nil {
		//	log.Errorf("error parse file:%v", err)
		//	return
		//} else {
		//	//http.ServeFile(w, r, "index.html")
		//	tmpl.Execute(w, map[string]interface{}{"uid": uid})
		//http.ServeFile(w, r, fmt.Sprintf("apps/im/home.html?user=%d", user))
		if user, userErr := userCli.Client.UserByToken(r.Context(), &shard.TokenReq{
			Token:              token,
		}); userErr != nil || user == nil {
			http.Error(w, "token not provide", http.StatusUnauthorized)
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			homeTempl.Execute(w, map[string]string{"host": r.Host, "token": token})
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
		if token := r.URL.Query().Get("token"); token == "" {
			http.Error(w, "token not provide", http.StatusUnauthorized)
		} else {
			//http.ServeFile(w, r, fmt.Sprintf("apps/im/home.html?user=%d", user))
			if user, userErr := userCli.Client.UserByToken(r.Context(), &shard.TokenReq{
				Token:              token,
			}); userErr != nil || user == nil {
				w.WriteHeader(400)
				w.Write([]byte("please login first"))
				return
			} else {
				log.Printf("ws connect:%d", user.UserId)
				serveWs(hub, w, r, user)
			}
		}
	})
	addr := fmt.Sprintf("0.0.0.0:%d", httpPort)
	log.Printf("listening addr:%s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
