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
	"net/http"
)

//ListenAndServe starts the listening process
func (hmc *HTTPMpc) ListenAndServe() {
	s := http.ListenAndServe(fmt.Sprintf(":%d", hmc.config.Port), hmc.router)
	fmt.Println(s)

}
func (hmc *HTTPMpc) hNext(w http.ResponseWriter, r *http.Request) {
	hmc.execute(w, r, hmc.mpd.Next)
}
func (hmc *HTTPMpc) hPrevious(w http.ResponseWriter, r *http.Request) {
	hmc.execute(w, r, hmc.mpd.Previous)
}
func (hmc *HTTPMpc) hPing(w http.ResponseWriter, r *http.Request) {
	hmc.execute(w, r, hmc.mpd.Ping)
}
func (hmc *HTTPMpc) hStop(w http.ResponseWriter, r *http.Request) {
	hmc.execute(w, r, hmc.mpd.Stop)
}

//The next few routes are boolean functions

func (hmc *HTTPMpc) hConsume(w http.ResponseWriter, r *http.Request) {
	hmc.setclear(w, r, hmc.mpd.Consume)
}
func (hmc *HTTPMpc) hPause(w http.ResponseWriter, r *http.Request) {
	hmc.setclear(w, r, hmc.mpd.Pause)
}
func (hmc *HTTPMpc) hRandom(w http.ResponseWriter, r *http.Request) {
	hmc.setclear(w, r, hmc.mpd.Random)
}
func (hmc *HTTPMpc) hRepeat(w http.ResponseWriter, r *http.Request) {
	hmc.setclear(w, r, hmc.mpd.Repeat)
}
func (hmc *HTTPMpc) hSingle(w http.ResponseWriter, r *http.Request) {
	hmc.setclear(w, r, hmc.mpd.Single)
}

//various functions that return attrs
func (hmc *HTTPMpc) hStatus(w http.ResponseWriter, r *http.Request) {
	hmc.attrs(w, r, hmc.mpd.Status)
}
func (hmc *HTTPMpc) hStats(w http.ResponseWriter, r *http.Request) {
	hmc.attrs(w, r, hmc.mpd.Stats)
}
func (hmc *HTTPMpc) hCurrentSong(w http.ResponseWriter, r *http.Request) {
	hmc.attrs(w, r, hmc.mpd.CurrentSong)
}

//attrsURIslice slice items
func (hmc *HTTPMpc) hFind(w http.ResponseWriter, r *http.Request) {
	hmc.attrsURISlice(w, r, hmc.mpd.Find)
}
func (hmc *HTTPMpc) hListInfo(w http.ResponseWriter, r *http.Request) {
	hmc.attrsURISlice(w, r, hmc.mpd.ListInfo)
}
func (hmc *HTTPMpc) hListAllInfo(w http.ResponseWriter, r *http.Request) {
	hmc.attrsURISlice(w, r, hmc.mpd.ListAllInfo)
}
func (hmc *HTTPMpc) hPlaylistContents(w http.ResponseWriter, r *http.Request) {
	hmc.attrsURISlice(w, r, hmc.mpd.PlaylistContents)
}

//attrsSlice
func (hmc *HTTPMpc) hListOutputs(w http.ResponseWriter, r *http.Request) {
	hmc.attrsSlice(w, r, hmc.mpd.ListOutputs)
}
func (hmc *HTTPMpc) hListPlaylists(w http.ResponseWriter, r *http.Request) {
	hmc.attrsSlice(w, r, hmc.mpd.ListPlaylists)
}
