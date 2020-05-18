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
package colors

import (
	"fmt"

	"github.com/fatih/color"
)

// Colors are a struct to define the different colors available for printing data
type Colors struct {
	colorInfo        *color.Color
	colorInfoHeading *color.Color
	colorOk          *color.Color
	colorError       *color.Color
}

// ColorsGlobal is the Global color struct. Need to be setup before usage
var ColorsGlobal Colors

// DisplayInfo display the str string as the info. Need to pass the color structure
func DisplayInfo(str string) {
	ColorsGlobal.colorInfo.Print(str)
	fmt.Println()
}

// DisplayInfoNoNL display the str string as the info without appending a new line. Need to pass the color structure
func DisplayInfoNoNL(str string) {
	ColorsGlobal.colorInfo.Print(str)
}

// DisplayInfoHeading display the str string as a heading info. Need to pass the color structure
func DisplayInfoHeading(str string) {
	ColorsGlobal.colorInfoHeading.Print(str)
	fmt.Println() // Because the newline keeps the background color, we need to do it separately
}

// DisplayError display the str string as an error. Need to pass the color structure
func DisplayError(str string) {
	ColorsGlobal.colorError.Print(str)
	fmt.Println()
}

// DisplayOk display the str string as a ok message. Need to pass the color structure
func DisplayOk(str string) {
	ColorsGlobal.colorOk.Println(str)
}

// DisplayOkNoNL display the str string as a ok message without appending a new line. Need to pass the color structure
func DisplayOkNoNL(str string) {
	ColorsGlobal.colorOk.Print(str)
}

// SetupColors set up the color struct given in parameter.
func SetupColors() {
	ColorsGlobal.colorInfo = color.New(color.FgBlue).Add(color.BgWhite)
	ColorsGlobal.colorInfoHeading = color.New(color.FgWhite).Add(color.BgBlue)
	ColorsGlobal.colorOk = color.New(color.FgGreen)
	ColorsGlobal.colorError = color.New(color.FgWhite).Add(color.BgHiRed)
}
