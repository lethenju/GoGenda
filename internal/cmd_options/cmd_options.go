/*
MIT License

Copyright (c) 2020 Julien LE THENO

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
/*
 ============= GOGENDA SOURCE CODE ===========
 @Description : GoGenda is a CLI for google agenda, to focus on one task at a time and logs your activity
 @Author : Julien LE THENO
 =============================================
*/
package cmdOptions

import (
	"errors"
	"flag"
)

var setOptions map[string]string

// Init parses the arguments with the 'flag' package and set them in the local private setOptions variable
func Init() []string {
	setOptions = make(map[string]string)

	shell := flag.Bool("i", false, "Interactive shell")
	help := flag.Bool("h", false, "Help")
	compact := flag.Bool("compact", false, "Compact output")
	config := flag.String("config", "", "Custom configuration")

	flag.Parse()

	if *shell {
		setOptions["shell"] = "true"
	}
	if *help {
		setOptions["help"] = "true"
	}
	if *compact {
		setOptions["compact"] = "true"
	}
	if *config != "" {
		setOptions["config"] = *config
	}
	return flag.Args()
}

// IsOptionSet checks if the option given in parameters had been set by the user
func IsOptionSet(option string) bool {
	return setOptions[option] != ""
}

// GetStringOption returns the string value of a option
func GetStringOption(option string) (string, error) {
	ret := setOptions[option]
	if ret == "" {
		return "", errors.New("This option is not set")
	}
	return ret, nil
}

// GetNumberOfOptions returns the number of options that had been set
func GetNumberOfOptions() int {
	return len(setOptions)
}
