/*
The MIT License (MIT)

Copyright (c) 2016 Nick Potts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package httpmpc

import (
	"encoding/json"
	"github.com/GeertJohan/go.rice"
	"github.com/fhs/gompd/mpd"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

//HTTPMpc is base class struct
type HTTPMpc struct {
	config configuration
	mpd    *mpd.Client
	router *mux.Router
	mutex  *sync.Mutex
}

//New returns a properly configured HttpMpc
func New(cfg configuration) (hmc *HTTPMpc, err error) {
	hmc = new(HTTPMpc)
	hmc.config = cfg
	router := mux.NewRouter()
	box := rice.MustFindBox("html")
	router.Handle("/", http.FileServer(box.HTTPBox()))
	router.HandleFunc("/next", hmc.hNext).Methods("GET")
	router.HandleFunc("/previous", hmc.hPrevious).Methods("GET")
	router.HandleFunc("/ping", hmc.hPing).Methods("GET")
	router.HandleFunc("/stop", hmc.hStop).Methods("GET")
	router.HandleFunc("/consume", hmc.hConsume).Methods("PUT", "DELETE")
	router.HandleFunc("/pause", hmc.hPause).Methods("PUT", "DELETE")
	router.HandleFunc("/random", hmc.hRandom).Methods("PUT", "DELETE")
	router.HandleFunc("/repeat", hmc.hRepeat).Methods("PUT", "DELETE")
	router.HandleFunc("/single", hmc.hSingle).Methods("PUT", "DELETE")
	router.HandleFunc("/status", hmc.hStatus).Methods("GET")
	router.HandleFunc("/stats", hmc.hStats).Methods("GET")
	router.HandleFunc("/currentsong", hmc.hCurrentSong).Methods("GET")
	hmc.router = router

	//Setup HTTP server
	hmc.mpd, err = mpd.DialAuthenticated("tcp", cfg.MpdDial, cfg.Password)
	return
}

func (hmc *HTTPMpc) execute(w http.ResponseWriter, r *http.Request, exec func() error) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	e := exec()
	if e == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (hmc *HTTPMpc) setclear(w http.ResponseWriter, r *http.Request, exec func(bool) error) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	var e error
	switch r.Method {
	case "PUT":
		e = exec(true)
	case "DELETE":
		e = exec(false)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if e == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (hmc *HTTPMpc) attrs(w http.ResponseWriter, r *http.Request, exec func() (mpd.Attrs, error)) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	if st, err := exec(); err == nil {
		b, err := json.Marshal(st)
		if err == nil {
			w.Write(b)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
