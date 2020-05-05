# GoGenda


A CLI application to quickly tell your google agenda what you're doing and let it logs your activity until you tell it to stop

In my opinion it is a good and easy way to stay productive by keeping on a task for a certain time

## Installation

TODO

## Usage

The CLI is really easy, just run gogenda and type help to know what you can do
```
~/go/src/github.com/lethenju/gogenda> go run gogenda.go 
Welcome to GoGenda!
> 
> 
> help
== GoGenda ==
 GoGenda helps you keep track of your activities
 = Commands = 
 START WORK - Start a work related activity
 START ORGA - Start a organisation related activity - 
                Reading articles, answering mails etc
 START LUNCH - Start a lunch related activity
 STOP - Stop the current activity
 RENAME - Rename the current activity
> start work
Enter name of event :  GoGenda Readme Redaction
Successfully added event ! Work hard! 
[GoGenda Readme Redaction]> 
[GoGenda Readme Redaction]> 
[GoGenda Readme Redaction]> rename
Enter name of event :  GoGenda : Readme redaction
Successfully renamed the activity
[GoGenda : Readme redaction]> 
[GoGenda : Readme redaction]> stop
Successfully stopped the activity ! I hope it went well 
> exit
See you later !
```
