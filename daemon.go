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
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/fhs/gompd/mpd"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//HTTPMpc is base class struct
type HTTPMpc struct {
	config configuration
	mpd    *mpd.Client
	router *mux.Router
	mutex  sync.Mutex
}

//New returns a properly configured HttpMpc
func New(cfg configuration) (hmc *HTTPMpc, err error) {
	hmc = new(HTTPMpc)
	hmc.config = cfg
	router := mux.NewRouter()
	router.HandleFunc("/next", hmc.hNext).Methods("POST")
	router.HandleFunc("/previous", hmc.hPrevious).Methods("POST")
	router.HandleFunc("/ping", hmc.hPing).Methods("GET")
	router.HandleFunc("/stop", hmc.hStop).Methods("POST")
	router.HandleFunc("/clear", hmc.hClear).Methods("POST")
	router.HandleFunc("/consume", hmc.hConsume).Methods("PUT", "DELETE")
	router.HandleFunc("/pause", hmc.hPause).Methods("PUT", "DELETE")
	router.HandleFunc("/random", hmc.hRandom).Methods("PUT", "DELETE")
	router.HandleFunc("/repeat", hmc.hRepeat).Methods("PUT", "DELETE")
	router.HandleFunc("/single", hmc.hSingle).Methods("PUT", "DELETE")
	router.HandleFunc("/status", hmc.hStatus).Methods("GET")
	router.HandleFunc("/stats", hmc.hStats).Methods("GET")
	router.HandleFunc("/currentsong", hmc.hCurrentSong).Methods("GET")
	//URI handlers
	router.HandleFunc("/find/{uri:.*}", hmc.hFind).Methods("GET")
	router.HandleFunc("/listinfo/{uri:.*}", hmc.hListInfo).Methods("GET")
	router.HandleFunc("/listallinfo/{uri:.*}", hmc.hListAllInfo).Methods("GET")
	router.HandleFunc("/playlistcontents/{uri:.*}", hmc.hPlaylistContents).Methods("GET")
	//other attr handlers
	router.HandleFunc("/listoutputs", hmc.hListOutputs).Methods("GET")
	router.HandleFunc("/listplaylists", hmc.hListPlaylists).Methods("GET")
	router.HandleFunc("/playlistinfo", hmc.hPlaylistInfo).Methods("GET")
	router.HandleFunc("/add/{uri:.*}", hmc.hAdd).Methods("POST")
	router.HandleFunc("/playlistclear/{uri:.*}", hmc.hPlaylistClear).Methods("POST")
	router.HandleFunc("/playlistremove/{uri:.*}", hmc.hPlaylistRemove).Methods("POST")
	router.HandleFunc("/playlistsave/{uri:.*}", hmc.hPlaylistSave).Methods("POST")
	// func(int) error handlers
	router.HandleFunc("/deleteid/{idpos:[-]{0,1}[0-9]*}", hmc.hDeleteID).Methods("POST")
	router.HandleFunc("/play/{idpos:[-]{0,1}[0-9]*}", hmc.hPlay).Methods("POST")
	router.HandleFunc("/playid/{idpos:[-]{0,1}[0-9]*}", hmc.hPlayID).Methods("POST")
	router.HandleFunc("/disableoutput/{idpos:[-]{0,1}[0-9]*}", hmc.hDisableOutput).Methods("POST")
	router.HandleFunc("/enableoutput/{idpos:[-]{0,1}[0-9]*}", hmc.hEnableOutput).Methods("POST")
	router.HandleFunc("/setvolume/{idpos:[-]{0,1}[0-9]*}", hmc.hSetVolume).Methods("POST")

	css := rice.MustFindBox("css")
	router.Handle("/css/{path:.*}", http.StripPrefix("/css/", http.FileServer(css.HTTPBox())))

	html := rice.MustFindBox("html")
	// router.Handle("/", http.StripPrefix("/html/", http.FileServer(html.HTTPBox())))
	router.Handle("/{path:.*}", http.FileServer(html.HTTPBox()))

	hmc.router = router

	//Setup HTTP server
	hmc.mpd, err = mpd.DialAuthenticated("tcp", hmc.config.MpdDial, hmc.config.Password)
	if err == nil {
		go hmc.busy()
	}
	return
}

func (hmc *HTTPMpc) busy() {
	//ping every 1 sec, attempting a redial if connection fails
	broken := false
	for {
		time.Sleep(time.Millisecond * time.Duration(hmc.config.KeepAlive))
		hmc.mutex.Lock()
		if e := hmc.mpd.Ping(); e != nil { //network error  attempt redial
			if !broken {
				broken = true
				fmt.Println("Connection lost")
			}
			if hmc.mpd, e = mpd.DialAuthenticated("tcp", hmc.config.MpdDial, hmc.config.Password); e == nil {
				fmt.Println("Connection re-established")
				broken = false
			}
		}
		hmc.mutex.Unlock()
	}
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
		b, err := json.MarshalIndent(st, "", "  ")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
func (hmc *HTTPMpc) attrsURISlice(w http.ResponseWriter, r *http.Request, exec func(string) ([]mpd.Attrs, error)) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	vars := mux.Vars(r)
	uri, ok := vars["uri"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if st, err := exec(uri); err == nil {
		b, err := json.MarshalIndent(st, "", "  ")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
func (hmc *HTTPMpc) uri(w http.ResponseWriter, r *http.Request, exec func(string) error) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	vars := mux.Vars(r)
	uri, ok := vars["uri"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := exec(uri); err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
func (hmc *HTTPMpc) attrsSlice(w http.ResponseWriter, r *http.Request, exec func() ([]mpd.Attrs, error)) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	if st, err := exec(); err == nil {
		b, err := json.MarshalIndent(st, "", "  ")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
func (hmc *HTTPMpc) int(w http.ResponseWriter, r *http.Request, exec func(int) error) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	vars := mux.Vars(r)
	if idpos, ok := vars["idpos"]; ok {
		if id, err := strconv.ParseInt(idpos, 10, 16); err == nil {
			if err := exec(int(id)); err == nil {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
func (hmc *HTTPMpc) hPlaylistInfo(w http.ResponseWriter, r *http.Request) {
	hmc.mutex.Lock()
	defer hmc.mutex.Unlock()
	start := int64(-1)
	end := int64(-1)
	// var err error
	// vars := mux.Vars(r)
	// if s, ok := vars["start"]; ok {
	// 	if start, err = strconv.ParseInt(s, 10, 32); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	// if s, ok := vars["end"]; ok {
	// 	if end, err = strconv.ParseInt(s, 10, 32); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	if st, err := hmc.mpd.PlaylistInfo(int(start), int(end)); err == nil {
		b, err := json.MarshalIndent(st, "", "  ")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
