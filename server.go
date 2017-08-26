package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	Token  string
	Status bool
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	if Token == "" {
		fmt.Println("No token. useage: server -t <bot token> ")
		return
	}

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error making the discord session: ", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection: ", err)
		return
	}

	fmt.Println("Bot is now running have fun!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!test" {
		s.ChannelMessageSend(m.ChannelID, "Fuck you "+m.Author.Username)
	}

	if m.Content == "!embed" {
		image := &discordgo.MessageEmbedImage{URL: "https://maps.googleapis.com/maps/api/staticmap?center=CollegeStation,TX&zoom=14&size=400x400"}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Title: "embeded message test", Color: 1000, Image: image})

	}

	if m.Content == "!track" {
		if Status == false {
			messageChannel(s, m.ChannelID, "Now trackning the aggie1 with aprs.fi")
			req, err := http.NewRequest("GET", "https://httpbin.org/ip", nil)

			req.Header.Add("user-agent", "APRSDiscordTracker/0.0.1")

			if err != nil {
				fmt.Println("there was an error :( :")
				fmt.Println(err)
			}

			response, _ := http.DefaultClient.Do(req)

			defer response.Body.Close()
			if response.StatusCode == 200 {
				bodyBytes, err2 := ioutil.ReadAll(response.Body)
				bodyString := string(bodyBytes)
				messageChannel(s, m.ChannelID, bodyString)
				if err2 != nil {
					fmt.Println("oh god")
				}
			}

			Status = true
		} else {
			messageChannel(s, m.ChannelID, "Already tracking balloon")
		}
		return
	}

	if m.Content == "!stoptrack" {
		messageChannel(s, m.ChannelID, "No longer tracking balloon")
		Status = false
	}

}
func messageChannel(s *discordgo.Session, channelID string, message string) {
	fmt.Println(message)
	s.ChannelMessageSend(channelID, message)
}
