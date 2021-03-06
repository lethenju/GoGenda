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
package utilities

import (
	"bufio"
	"fmt"
	"os"
)

// InputFromUser is a helper function to ask nicely the user of some string to enter and get it
func InputFromUser(name string) (inputUser string) {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter " + name + " :")
	if !scanner.Scan() {
		return
	}
	return scanner.Text()
}

// AskOkFromUser is a helper function to ask nicely the user if he/she's okay to perform some action
func AskOkFromUser(str string) bool {

	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	for answer != "y" && answer != "n" {
		fmt.Print(str + " (y/n) :")
		scanner.Scan()
		answer = scanner.Text()
	}
	return answer == "y"
}
