package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var Caobj cache = NewCache()

//主节点对外开放的读写接口
func ListenServer() {

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		val := Caobj.Get(string(key))

		Jsn, err := json.Marshal(val)
		if err != nil {
			w.Write([]byte(Jsn))
		}
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		timOut, err := strconv.Atoi(r.PostForm.Get("timeout"))
		if err != nil {
			fmt.Println(err.Error())
		}
		key := r.PostForm.Get("key")

		Caobj.mu.Lock()
		Caobj.Put(key, r.PostForm.Get("data"), int64(timOut))
		Caobj.mu.Unlock()
	})

	http.ListenAndServe(":9991", nil)
}
