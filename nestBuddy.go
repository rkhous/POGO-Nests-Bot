package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

//CONFIG
var Token string // edit the bot token on the first line under the main function (func main())
var Prefix string = "-" //Prefix for the bot. Limited to one char.
var nestChannel string = "7246846588445897" //Create a nest channel on discord, right click, and copy the ID.
var listOfNestID = []string{} // DO NOT TOUCH THIS! Leave blank
var admins = []string{"ADMIN1297369", "ADMIN92736947"} //Put a list of users ID's who can use the migrate feature here. Read the readme for formatting.
var nestLocations = map[string]string{
  "Lake of Elves":"43.123123,132.123123",
  "Lake of Dragons":"43.456456,132.456456",
} // set up your nests and locations here, names should be lowercase, no space between lat,lon!

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	var Token string = "BOT TOKEN HERE" // Bot Token
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) == 0 {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	getArgs := strings.Split(m.Content, " ")
	if checkIfInCommand(getArgs[0]) && getArgs[0] == Prefix+"addnest" {
		if len(getArgs) >= 3 {
			pokemon := getArgs[1]
			locationName := strings.Split(m.Message.Content, pokemon+" ")[1]
			if doesNestExist(locationName, nestLocations) {
				locationCoords := getURLLocation(locationName)
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Color:  0x00ff00,
					Description: "**Reported by:** " + m.Author.Mention() + "\n\n" +
						"**Nesting Pokémon:** " + strings.Title(pokemon),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://img.pokemondb.net/sprites/x-y/normal/" + strings.ToLower(pokemon) + ".png",
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Created by github.com/rkhous",
						IconURL: "https://d1q6f0aelx0por.cloudfront.net/product-logos/81630ec2-d253-4eb2-b36c-eb54072cb8d6-golang.png"},
					Title: "**" + strings.Title(locationName) + "**",
					URL:   locationCoords,
				}
				sendNestMessage, err := s.ChannelMessageSendEmbed(nestChannel, embed)
				if err != nil {
					fmt.Println(err.Error())
					return
				} else {
					listOfNestID = append(listOfNestID, sendNestMessage.ID)
					s.ChannelMessageSend(m.ChannelID, "Okay, "+m.Author.Mention()+", I have successfully added "+
						"your report to the nests channel. Thank you! :blush:")
					fmt.Println("Nest reported: " + pokemon +
						"\nReported by: " + m.Author.Username + "\n")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Sorry, "+m.Author.Mention()+", but `"+
					strings.Title(locationName)+"` is not a nesting area I'm aware of.")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "An incorrect number of arguments were given.")
		}
	} else if checkIfInCommand(getArgs[0]) && getArgs[0] == Prefix+"undo" {
		if len(getArgs) >= 3 {
			s.ChannelMessageSend(m.ChannelID, "Incorrect number of arguments given.\n"+
				"Please use the following format: `@undo <message_id>`")
		} else {
			msgID := getArgs[1]
			indexOfMsg := findMessage(msgID, listOfNestID)
			if indexOfMsg == -1 {
				s.ChannelMessageSend(m.ChannelID, "Sorry! But I couldn't find that message in the current list of nest spawns.")
			} else {
				s.ChannelMessageDelete(nestChannel, listOfNestID[indexOfMsg])
				listOfNestID = append(listOfNestID[:indexOfMsg], listOfNestID[indexOfMsg+1:]...)
				s.ChannelMessageSend(m.ChannelID, "Successfully deleted nest data from list and removed message from nest channel.")
				fmt.Println("ID Removed: " + msgID + "\nRemoved by: " + m.Message.Author.Username)
			}
		}
	} else if checkIfInCommand(getArgs[0]) && getArgs[0] == Prefix+"migrate" {
		if isAdmin(m.Author.ID, admins) {
			s.ChannelMessagesBulkDelete(nestChannel, listOfNestID)
			listOfNestID = []string{}
			s.ChannelMessageSend(m.ChannelID, "Okay, "+m.Author.Mention()+"! I have deleted all the old nests from the nest channel!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Sorry, but this command is only limited to a few users.")
		}
	} else if checkIfInCommand(getArgs[0]) && getArgs[0] == Prefix+"help"{
		help_embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0x00ff00,
			Description: "**Commands and Examples:**\n\n" +
				"**Adding a Nest**\n" +
				"Command: `" + Prefix + "addnest <pokemon> <location name of nest>`\n" +
				"Example: " + Prefix + "addnest pikachu griffith park\n\n" +
				"**Handling Mistakes**\n" +
				"Command: `" + Prefix + "undo <message id>`\n" +
				"Example: " + Prefix + "undo 518231154097127444\n\n" +
				"**Nest Migrations - Limited to Staff**\n" +
				"Command: `" + Prefix + "migrate`\n\n" +
				"**See List of Nests**\n" +
				"Command: `" + Prefix + "list`",
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Created by github.com/rkhous",
				IconURL: "https://d1q6f0aelx0por.cloudfront.net/product-logos/81630ec2-d253-4eb2-b36c-eb54072cb8d6-golang.png"},
			Title: "**Nest Report Help**",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, help_embed)
	}else if checkIfInCommand(getArgs[0]) && getArgs[0] == Prefix+"list"{
		s.ChannelMessageSend(m.ChannelID, mapToString(nestLocations))
	}else {
		return
	}
}

func checkIfInCommand(command string) bool {
	var listOfCommands = [] string {"addnest", "undo", "migrate", "help", "list"}
	for _, n := range listOfCommands{
		if strings.ToLower(Prefix + n) == strings.ToLower(command){
			return true
		}else{
			continue
		}
	}
	return false
}

func findMessage(messageID string, a [] string) int {
	for i := 0; i < len(a); i++{
		if a[i] == messageID{
			return i
		}else{
			continue
		}
	}
	return -1
}

func getURLLocation(name string) string{
	name = strings.ToLower(name)
	return ("http://maps.google.com/maps?daddr=" + nestLocations[name] + "&amp;ll=")
}

func doesNestExist(name string, m map[string]string) bool {
	name = strings.ToLower(name)
	for n := range m{
		if n == name{
			return true
		}else{
			continue
		}
	}
	return false
}

func isAdmin(userID string, ls[]string) bool {
	for _, n := range ls{
		if n == userID{
			return true
		}else{
			continue
		}
	}
	return false
}

func mapToString(m map[string]string) string{
	mapKeys := []string{}
	for n := range m {
		mapKeys = append(mapKeys, n)
	}
	return strings.Title(strings.Join(mapKeys, ",\n"))
}
