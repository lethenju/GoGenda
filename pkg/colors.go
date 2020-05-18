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
package main

import (
	"fmt"

	"github.com/fatih/color"
)

type colors struct {
	colorInfo        *color.Color
	colorInfoHeading *color.Color
	colorOk          *color.Color
	colorError       *color.Color
}

func displayInfo(colors *colors, str string) {
	colors.colorInfo.Print(str)
	fmt.Println()
}
func displayInfoNoNL(colors *colors, str string) {
	colors.colorInfo.Print(str)
}
func displayInfoHeading(colors *colors, str string) {
	colors.colorInfoHeading.Print(str)
	fmt.Println() // Because the newline keeps the background color, we need to do it separately
}
func displayError(colors *colors, str string) {
	colors.colorError.Print(str)
	fmt.Println()
}
func displayOk(colors *colors, str string) {
	colors.colorOk.Println(str)
}
func displayOkNoNL(colors *colors, str string) {
	colors.colorOk.Print(str)
}

func setupColors(colors *colors) {
	colors.colorInfo = color.New(color.FgBlue).Add(color.BgWhite)
	colors.colorInfoHeading = color.New(color.FgWhite).Add(color.BgBlue)
	colors.colorOk = color.New(color.FgGreen)
	colors.colorError = color.New(color.FgWhite).Add(color.BgHiRed)
}
