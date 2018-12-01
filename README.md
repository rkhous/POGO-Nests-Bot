# POGO Nests Bot
A bot, written in Go, to help discord servers efficiently, neatly, and quickly access and set up nests to share with their users.

## Setup
- Install GoLang: https://golang.org/doc/install
- Create a folder wherever you plan on running the bot, and download the files into that folder.
- Configure your $GOPATH - Do this in the folder you install the file
- If you're having trouble installing, please look up the issue on Stackoverflow.
- Next, install DiscordGO by running this in terminal: go get github.com/bwmarrin/discordgo

## Configuring the Bot
- There's a few things you will have to edit inside nestBuddy.go
- First, get your discord bots token and replace "Bot Token Here" with your bots token, under the main function. 
- Next, create a channel where the bot will be posting the nests, and copy it's ID under the //CONFIG comment.
- Next, a list of admins must be added to be able to use the migrate function. Use the template below to follow format. 
- _var admins = []string{"Admin ID 1", "Admin ID 2", "etc"}_

## Adding Nests to Code
- Finally, you will need to add your area's nests so the bot can grab locations. Please edit this portion as shown below:
```Go
var nestLocations = map[string]string{
	"Park of Dragons":"43.00009,-121.744121",
	"Lake of Elves":"43.152435,-121.260806",
}
```

## Developer
- Robert Khousadian
- If you'd like to buy me a coffee: https://paypal.me/rkhous
- If you have questions or a problem, join my discord or open an issue: https://discord.gg/CwqbHt5
