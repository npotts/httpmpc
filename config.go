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
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type configuration struct {
	MpdDial   string `yaml:"MPD"`
	Password  string `yaml:"Password"`
	Port      int    `yaml:"HTTP Port"`
	KeepAlive int    `yaml:"Keep Alive"`
}

//BaseConfig is the preloaded config
var BaseConfig = configuration{MpdDial: "localhost:6600", Password: "", Port: 8080}

//if set to a non-empty string, will read from this config file
var thisfile string

func configHunt() []byte {
	if thisfile != "" {
		bytes, err := ioutil.ReadFile(thisfile)
		if err == nil {
			return bytes
		}
		panic(fmt.Errorf("Unable to read %q: %v", thisfile, err))
	}
	for _, file := range []string{"", "$HOME/.", "/etc/"} {
		file = os.ExpandEnv(file) + "httpmpc.yml"
		if bytes, err := ioutil.ReadFile(file); err == nil {
			fmt.Printf("Using config file %q\n", file)
			return bytes
		}
	}
	panic(fmt.Errorf("Unable to locate %q", "httpmpc.yml"))
}

func init() {
	flag.StringVar(&thisfile, "config", "", "If specified, use this config file")
}

//Get attempts to read a config from somewhere
func Get() (c configuration) {
	if e := yaml.Unmarshal(configHunt(), &c); e != nil {
		panic(e)
	}
	return
}
