# GoGenda


A CLI application to quickly tell your google agenda what you're doing and let it logs your activity until you tell it to stop

In my opinion it is a good and easy way to stay productive by keeping on a task for a certain time

## Installation

TODO

## Usage

Create a configuration json  named `config.json` in `~/.gogenda/` folder. It will keep track of the categories of activities you want to log, and the colors you want to have in your google agenda

My example : 
```json
{
    "categories":
    [
        {
            "name":"WORK",
            "color":"red"
        },
        {
            "name":"ORGA",
            "color":"yellow"
        },
        {
            "name":"LUNCH",
            "color":"purple"
        },
        {
            "name":"FUN",
            "color":"orange"
        }
    ]
}
```


The CLI is really easy, just run gogenda and type help to know what you can do
(it actually has color, this is just a copy paste)
```
~ $ gogenda
Welcome to GoGenda!
Version number : 0.1.4
> 
> 
> help
== GoGenda ==
 GoGenda helps you keep track of your activities
 = Commands = 

 START WORK - Add an event in red
 START ORGA - Add an event in yellow
 START LUNCH - Add an event in purple
 START FUN - Add an event in orange
 STOP - Stop the current activity
 RENAME - Rename the current activity
 DELETE - Delete the current activity
> 
> start work this is my work
Successfully added activity ! 
[ this is my work 1s ]> 
[ this is my work 3s ]> 
[ this is my work 4s ]> stop
The activity 'this is my work' lasted 9s
Successfully stopped the activity ! I hope it went well 
> exit
See you later !
~ $ 
```
