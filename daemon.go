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
	"fmt"
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
	//Add Routes
	// router.HandleFunc(path, f)
	hmc.router = router
	//Setup HTTP server
	hmc.mpd, err = mpd.DialAuthenticated("tcp", cfg.MpdDial, cfg.Password)
	return
}

//ListenAndServe starts the listening process
func (hmc *HTTPMpc) ListenAndServe() {
	http.ListenAndServe(fmt.Sprintf(":%d", hmc.config.Port), hmc.router)
}
