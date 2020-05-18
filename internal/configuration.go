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
package gogenda

import (
	"encoding/json"
	"os"
	"strings"
)

// ConfigCategory is a category of activity
type ConfigCategory struct {

	// Name
	Name string `json:"name"`
	// Color
	Color string `json:"color"`
}

// Config represents the configuration of the app
type Config struct {
	// Categories are the active categories of activities
	Categories []ConfigCategory `json:"categories"`
}

// Conf is the globally accessible configuration
var Conf Config

// LoadConfiguration : Init and load the configuration in the conf variable
func LoadConfiguration(file string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	Conf = Config{}
	err = json.NewDecoder(f).Decode(Conf)
	return err
}

//GetColorFromName returns the color string for a name category
func GetColorFromName(name string) (color string) {

	ourCategory := ConfigCategory{Name: "default", Color: "blue"}
	for _, category := range Conf.Categories {
		if strings.ToUpper(name) == category.Name {
			ourCategory = category
		}
	}
	return ourCategory.Color
}

//GetNameFromColor returns the name string for a color category
func GetNameFromColor(color string) (name string) {
	ourCategory := ConfigCategory{Name: "default", Color: "blue"}
	for _, category := range Conf.Categories {
		if color == category.Color {
			ourCategory = category
		}
	}
	return ourCategory.Name
}
