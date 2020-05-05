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
 @Version : 0.1.4
 @Author : Julien LE THENO
 =============================================
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

func displayInfo(ctx *gogendaContext, str string) {
	ctx.colors.colorInfo.Print(str)
	fmt.Println()
}
func displayInfoNoNL(ctx *gogendaContext, str string) {
	ctx.colors.colorInfo.Print(str)
}
func displayInfoHeading(ctx *gogendaContext, str string) {
	ctx.colors.colorInfoHeading.Print(str)
	fmt.Println() // Because the newline keeps the background color, we need to do it separately
}
func displayError(ctx *gogendaContext, str string) {
	ctx.colors.colorError.Print(str)
	fmt.Println()
}
func displayOk(ctx *gogendaContext, str string) {
	ctx.colors.colorOk.Println(str)
}
func displayOkNoNL(ctx *gogendaContext, str string) {
	ctx.colors.colorOk.Print(str)
}

func setupColors(ctx *gogendaContext) {
	ctx.colors.colorInfo = color.New(color.FgBlue).Add(color.BgWhite)
	ctx.colors.colorInfoHeading = color.New(color.FgWhite).Add(color.BgBlue)
	ctx.colors.colorOk = color.New(color.FgGreen)
	ctx.colors.colorError = color.New(color.FgWhite).Add(color.BgHiRed)
}
