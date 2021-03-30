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
func (hmc *HTTPMpc) hClear(w http.ResponseWriter, r *http.Request) {
	hmc.execute(w, r, hmc.mpd.Clear)
}

//The next few routes are boolean functions, that use a CommandList
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
	cmdlist := hmc.mpd.BeginCommandList()
	cmdlist.Single(hmc.methodToBool(r))
	if err := cmdlist.End(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (hmc *HTTPMpc) hConsume(w http.ResponseWriter, r *http.Request) {
	cmdlist := hmc.mpd.BeginCommandList()
	cmdlist.Consume(hmc.methodToBool(r))
	if err := cmdlist.End(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	//hmc.attrsURISlice(w, r, hmc.mpd.Find)
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

func (hmc *HTTPMpc) hAdd(w http.ResponseWriter, r *http.Request) {
	hmc.uri(w, r, hmc.mpd.Add)
}
func (hmc *HTTPMpc) hPlaylistClear(w http.ResponseWriter, r *http.Request) {
	hmc.uri(w, r, hmc.mpd.PlaylistClear)
}
func (hmc *HTTPMpc) hPlaylistRemove(w http.ResponseWriter, r *http.Request) {
	hmc.uri(w, r, hmc.mpd.PlaylistRemove)
}

func (hmc *HTTPMpc) hPlaylistSave(w http.ResponseWriter, r *http.Request) {
	hmc.uri(w, r, hmc.mpd.PlaylistSave)
}

//func(int)error responses
func (hmc *HTTPMpc) hDeleteID(w http.ResponseWriter, r *http.Request) {
	hmc.int(w, r, hmc.mpd.DeleteID)
}
func (hmc *HTTPMpc) hPlay(w http.ResponseWriter, r *http.Request)   { hmc.int(w, r, hmc.mpd.Play) }
func (hmc *HTTPMpc) hPlayID(w http.ResponseWriter, r *http.Request) { hmc.int(w, r, hmc.mpd.PlayID) }
func (hmc *HTTPMpc) hDisableOutput(w http.ResponseWriter, r *http.Request) {
	hmc.int(w, r, hmc.mpd.DisableOutput)
}
func (hmc *HTTPMpc) hEnableOutput(w http.ResponseWriter, r *http.Request) {
	hmc.int(w, r, hmc.mpd.EnableOutput)
}
func (hmc *HTTPMpc) hSetVolume(w http.ResponseWriter, r *http.Request) {
	hmc.int(w, r, hmc.mpd.SetVolume)
}
