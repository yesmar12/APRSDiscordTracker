package main

import (
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	Token   string
	Status  bool
	Station string
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

	if m.Content == "!help" {
		msg := "There are 5 commands for this bot. [!help, !embed, !starttrack, !track, !stoptrack]\n !starttrack: starts the tracking system currently starting tracking does nothing\n!track: prints the map of the location of the bot\n!stoptrack: stops tracking"
		s.ChannelMessageSend(m.ChannelID, "Hi "+m.Author.Username)
		s.ChannelMessageSend(m.ChannelID, msg)

	}

	if m.Content == "!embed" {
		image := &discordgo.MessageEmbedImage{URL: "https://maps.googleapis.com/maps/api/staticmap?center=30.61950,-96.33933&zoom=15&size=1000x1000&maptype=satellite&key=AIzaSyDLjIV_Io8-QWuMbbRnkR3GQvvtZcFdGZY&markers=blue|30.61950,-96.33933"}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Title: "embeded message test", Color: 1000, Image: image})
	}

	if strings.HasPrefix(m.Content, "!starttrack") {
		if Status == false {
			Station = strings.Replace(strings.Replace(m.Content, "!starttrack", "", -1), " ", "", -1)
			if Station == "" {
				Station = "AGGIE-2"
				messageChannel(s, m.ChannelID, "Using default APRS Value:\" "+Station+"\"")
			}
			messageChannel(s, m.ChannelID, "Tracking APRS station: \""+Station+"\"")
			req, err := http.NewRequest("GET", "https://api.aprs.fi/api/get?name="+Station+"&what=loc&apikey=102675.0ZXtEN0HRMaAl&format=json", nil)

			req.Header.Add("user-agent", "APRSDiscordTracker/0.0.1")

			if err != nil {
				fmt.Println("there was an error :( :")
				fmt.Println(err)
			}

			response, _ := http.DefaultClient.Do(req)

			defer response.Body.Close()
			if response.StatusCode == 200 {
				bodyBytes, err2 := ioutil.ReadAll(response.Body)
				lat, err3 := jsonparser.GetString(bodyBytes, "entries", "[0]", "lat")
				lng, err4 := jsonparser.GetString(bodyBytes, "entries", "[0]", "lng")
				msg := "This is the lat and lng:" + lat + "," + lng
				messageChannel(s, m.ChannelID, msg)
				if err2 != nil || err3 != nil || err4 != nil {
					fmt.Println("oh god")
				}
			}

			Status = true
		} else {
			messageChannel(s, m.ChannelID, "Already tracking balloon")
		}
		return

	}

	if m.Content == "!track" {
		if Status == true {
			req, err := http.NewRequest("GET", "https://api.aprs.fi/api/get?name="+Station+"&what=loc&apikey=102675.0ZXtEN0HRMaAl&format=json", nil)

			req.Header.Add("user-agent", "APRSDiscordTracker/0.0.1")

			if err != nil {
				fmt.Println("there was an error :( :")
				fmt.Println(err)
			}

			response, _ := http.DefaultClient.Do(req)

			defer response.Body.Close()
			if response.StatusCode == 200 {
				fmt.Println("Response code good")
				bodyBytes, err2 := ioutil.ReadAll(response.Body)
				bodyString := string(bodyBytes)
				fmt.Println(bodyString)
				lat, err3 := jsonparser.GetString(bodyBytes, "entries", "[0]", "lat")
				lng, err4 := jsonparser.GetString(bodyBytes, "entries", "[0]", "lng")
				name, err5 := jsonparser.GetString(bodyBytes, "entries", "[0]", "name")
				image := &discordgo.MessageEmbedImage{URL: "https://maps.googleapis.com/maps/api/staticmap?center=" + lat + "," + lng + "&zoom=15&size=1000x1000&maptype=satellite&key=AIzaSyDLjIV_Io8-QWuMbbRnkR3GQvvtZcFdGZY&markers=blue|" + lat + "," + lng}
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{URL: "https://www.google.com/maps/search/48.96550,2.24217", Title: "Current location of " + name + "(Click this for google maps directions)", Color: 1000, Image: image})
				if err2 != nil || err3 != nil || err4 != nil || err5 != nil {
					fmt.Println("oh god")
				}
			} else {
				fmt.Println("response code:", response.StatusCode)
			}

		} else {
			messageChannel(s, m.ChannelID, "Use !starttrack to begin!")
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
