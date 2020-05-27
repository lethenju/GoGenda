![](https://lethenju.github.io/res/logo_gogenda.png)

# GoGenda

A CLI application to quickly tell your google agenda what you're doing and let it logs your activity until you tell it to stop

In my opinion it is a good and easy way to stay productive by keeping on a task for a certain time

On top of that first use case, Gogenda now can give you statistics about the time you spent.
You can also manipulate any event of your calendar, move it from a date/time to another, renaming it and deleting it.

That way you can automate your calendar in much more advanced ways that the Google Agenda web interface.

## Installation

Install the app with 
```sh
git clone git@github.com:lethenju/gogenda
mkdir "~/.gogenda/"
cd cmd/gogenda
sudo GOBIN=/usr/local/bin go install
gogenda
```

GoGenda uses the REST Api of Google Calendar.

For obvious security reasons I cannot give you the keys of the application to log in the API.

But you can register your own application [here](https://console.developers.google.com/apis/credentials/wizard?)

Once you have the `credentials.json` file, put it in `.gogenda/` and launch gogenda.
It will ask you to put a link in your browser to allow your app to connect to your google account.

Then normally everything should work :) 

## Usage

### Configuration
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

### CLI Presentation

The CLI is really easy, just run gogenda for help
(it actually has color, this is just a copy paste)
```
$: gogenda
== GoGenda ==
 GoGenda helps you keep track of your activities
 Type Gogenda -h (command) to have more help for a specific command

 = Options = 
Important : options have to be used before command arguments !
 gogenda -i              - Launch the shell UI
 gogenda -h              - shows the help
 gogenda -compact        - Have minimalist output
 gogenda -config='path'  - Use a custom config file (absolute path only)

 = Commands = 
 gogenda start WORK - Add an event in red
 gogenda start ORGA - Add an event in yellow
 gogenda start LUNCH - Add an event in purple
 gogenda start FUN - Add an event in orange
 gogenda stop - Stop the current activity
 gogenda rename - Rename the current activity
 gogenda delete - Delete the current activity
 gogenda plan - See and manipulate your calendar as you want
 gogenda stats - shows statistics about your time spent in each category
 gogenda add - add an event to the planning. You can call it alone or with some params.
 gogenda help - show gogenda help (add a command name if you want specific command help)
```

### Shell Mode

You can use gogenda directly or launch an integrated shell (that way you have more info about the current launched event and you dont have to type gogenda each time you want to launch a command)


```
$: gogenda -i
Welcome to GoGenda!
Version number : 0.3.0
Last event : gogenda new readme
Are you still doing that ? (y/n) :y
[ gogenda new readme 7m41s ]> 
[ gogenda new readme 7m43s ]> 
```

### Gogenda Plan

The command `gogenda plan` gives you the ability to modify your calendar as you wish.
You have 4 sub-commands : `show`, `rename`, `move` and `delete` 
You have to call `gogenda plan show (your date)` to have the ID of the event you want to modify. You cannot modify an event with just its name or date, as several events can be under that description.

Type `gogenda help plan` to have more information about how to use it.

If you want to add an event to a custom date, use `gogenda add`.

If you want to add an event starting from now, use `gogenda start` instead.

### Gogenda Stats

You can also have some statistics about the time you spent on each category of your work for a given period.

For example, for me, yesterday :
```sh
$: gogenda stats yesterday
=== WORK ===
 [ 23:25 -> 00:01 ] 35m53s : opengl_framework debug
 [ 19:58 -> 21:50 ] 1h52m47s : plan commands
 [ 15:15 -> 15:34 ] 19m49s : OpenGL Doc
 [ 15:52 -> 16:09 ] 16m58s : siteperso
 [ 16:41 -> 17:07 ] 25m33s : messagingApp
 [ 17:07 -> 17:17 ] 9m27s : gogenda fix stats format
 [ 23:29 -> 23:31 ] 1m55s : fix stop
      Total : 3h42m22s
=== FUN ===
 [ 22:13 -> 22:26 ] 13m0s : youtube
      Total : 13m0s
```

You can also have a compact version, without the listed events : 

```sh
$: gogenda -compact stats yesterday
=== WORK ===
      Total : 3h42m22s
=== FUN ===
      Total : 13m0s
```