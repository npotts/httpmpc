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
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type configuration struct {
	MpdDial   string `yaml:"MPD"`
	Password  string `yaml:"Password"`
	Port      int    `yaml:"HTTP Port"`
	KeepAlive int    `yaml:"Keep Alive"`
}

//base is the preloaded config
var base = configuration{MpdDial: "localhost:6600", Password: "", Port: 8080, KeepAlive: 1000}

//if set to a non-empty string, will read from this config file
var thisfile string

func getConfig() (c configuration) {
	c = base
	bytes := []byte{}
	var err error

	if thisfile != "" {
		if bytes, err = ioutil.ReadFile(thisfile); err != nil {
			panic(fmt.Errorf("Unable to read %q: %v", thisfile, err))
		}
		if err = yaml.Unmarshal(bytes, &c); err != nil {
			panic(fmt.Errorf("Unable to parse %q: %v", thisfile, err))
		}
		return
	}

	for _, file := range []string{"", "$HOME/.", "/etc/"} {
		file = os.ExpandEnv(file) + "httpmpc.yml"
		if bytes, err = ioutil.ReadFile(file); err != nil {
			continue
		}
		if err := yaml.Unmarshal(bytes, &c); err != nil {
			fmt.Printf("Unable to parse %q:%v - reverting to default config", file, err)
			c = base //bad YAML use default
			continue
		} else {
			return //skip file
		}
	}
	return
}

var dflt string

func init() {
	flag.StringVar(&thisfile, "config", "", "If specified, use this config file")
	flag.StringVar(&dflt, "default", "", "If specified, writes a default configuration file to the specified file")
	flag.StringVar(&base.MpdDial, "mpd", "localhost:6600", "Connect to this mpd instance")
	flag.StringVar(&base.Password, "password", "", "Use this password.  If blank, it will be ignored")
	flag.IntVar(&base.Port, "http", 8080, "Serve out the HTTP out on this port")
	flag.IntVar(&base.KeepAlive, "keepalive", 1000, "Ensure communication with the mpd server at least this often")
}

var cfg = `--- 

#MPD instance, in server:port form.  Standard port is 6600
MPD: 192.168.42.74:6600

#Password, if required, to connect
Password:

#Which HTTP port to listen on
HTTP Port: 8080

#Keep Alive is how long we should wait (in ms)  when polling the MPD server to keep the connection Alive
Keep Alive: 1000
`

//Parse parses CLI flags and starts the daemon
func Parse() {
	flag.Parse()

	if dflt != "" { //write out default configuration and exit
		fo, err := os.Create(dflt)
		if err != nil {
			panic(err)
		}
		fo.WriteString(cfg)
		fo.Close()
		return
	}

	hmc, err := New(getConfig())
	fmt.Printf("httpmpc mpd:%q KeepAlive=%dms. HTTP Port:%d\n", hmc.config.MpdDial, hmc.config.KeepAlive, hmc.config.Port)
	if err != nil {
		panic(err)
	}

	hmc.ListenAndServe()
}
