package main

import (
	"log"
	"os"
	//"os/signal"
	//"syscall"
	"strings"
	"strconv"
	"bufio"
	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"slices"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	discord "github.com/bwmarrin/discordgo"
)

type Config struct {
	Token string
	Words []string
	Replacements []struct { 
		ToReplace []string
		Replacement string
	}
}

var (
	logger *log.Logger
	config Config
	x, y int
)

func main() {
	x, y = 25, 25
	file, err := os.OpenFile("message.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	logger = log.New(file, "", 0)

	cs, err := ioutil.ReadFile("config.json")
	eror(err)
	eror(json.Unmarshal([]byte(cs), &config))	

	botuser := "a"
	dg := newDiscordSession(config.Token)
	dm, err := dg.UserChannelCreate(botuser)
	eror(err)
	newDiscordAi(botuser, dm.ID, dg)

	// Wait here until CTRL-C or other term signal is received.
	//fmt.Println("Bot is now running. Press CTRL-C to exit.")
	//sc := make(chan os.Signal, 1)
	//signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	//<-sc

	reader := bufio.NewReader(os.Stdin)

	lastcall := ""
	for {
		text := strings.Replace(readMultiLineInput(reader), "\n", "", -1)

		if text != "!" {
			lastcall = text
		} else {
			text = lastcall
			fmt.Print(">>> "+text+"\n")
		}

		if strings.Contains(text, "!wordx ") {
			x, _ = strconv.Atoi(strings.Replace(text, "!wordx ", "", -1))
		} else if strings.Contains(text, "!wordy ") {
			y, _ = strconv.Atoi(strings.Replace(text, "!wordy ", "", -1))
		} else if strings.Contains(text, "!word ") {
			board := generateBoard(x, y, strings.Replace(text, "!word ", "", -1))

			if strings.Contains(board, "word cannot fit with board size"){
				fmt.Print("SY: error sending word search\n")
			} else {
				dg.ChannelMessageSend(dm.ID, "```"+board+"```")
			}
		} else {
			dg.ChannelMessageSend(dm.ID, text)

			bytes := make([]byte, len(text))
			normalize := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
				return unicode.Is(unicode.Mn, r)
			}), norm.NFC)
			_, _, err = normalize.Transform(bytes, []byte(text), true)
			eror(err)
			normalizedText := string(bytes)

			var words = []string{}

			for _, word := range config.Words {
				filterText := ""
				lastChar := ' '
				for _, char := range word {
					if lastChar != char{
						filterText+=string(char)
						lastChar = char
					}
				}

				replacedText := filterText
				for _, pair := range config.Replacements {
					for _, toReplace := range pair.ToReplace {
						replacedText = strings.ReplaceAll(replacedText, toReplace, pair.Replacement)
					}
				}

				filterText = ""
				lastChar = ' '
				for _, char := range replacedText {
					if lastChar != char{
						filterText+=string(char)
						lastChar = char
					}
				}
				words = append(words, filterText)
			}

			finalText := ""
			lastChar := ' '
			for _, char := range normalizedText {
				if lastChar != char{
					finalText+=string(char)
					lastChar = char
				}
			}

			replacedText := finalText
			for _, pair := range config.Replacements {
				for _, toReplace := range pair.ToReplace {
					replacedText = strings.ReplaceAll(replacedText, toReplace, pair.Replacement)
				}
			} 

			var caught string
			for _, word := range words {
				finalerText := ""
				var runes []rune

				for _, char := range word {
					runes = append(runes, char)
				}

				for _, char := range replacedText {
					if slices.Contains(runes, char){
						finalerText += string(char)
					}
				}

				finalestText := ""
				lastChar = ' '
				for _, char := range finalerText {
					if lastChar != char{
						finalestText+=string(char)
						lastChar = char
					}
				}

				if strings.Contains(finalestText, word) {
					caught = word
				}
			}

			if len(caught) > 0 {
				fmt.Print("FILTERED\n")
				//logger.Println(m.Author.Username+": "+m.Content)
				//logger.Println("caught: ", caught)
				//s.ChannelMessageDelete(m.ChannelID, m.ID)
			}
		}
		
	}

	// Cleanly close down the Discord session.
	dg.Close()
}

func readMultiLineInput(reader *bufio.Reader) string {
	fmt.Print(">>> ")

	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			os.Exit(0)
		}
		fmt.Printf("Reading the prompt failed: %s", err)
		os.Exit(1)
	}

	return line
}

func newDiscordSession(bottoken string) *discord.Session {
	dg, err := discord.New("Bot " + bottoken)
	eror(err)
	dg.Identify.Intents = discord.IntentsGuildMessages | discord.IntentDirectMessages
	eror(dg.Open())
	return dg
}

func newDiscordAi(botuser, dmid string, dg *discord.Session,){
	message := "a"
	//dg.ChannelMessageSend(dmid, message)
	dm, err := dg.UserChannelCreate(botuser)
	eror(err)
	newDiscordMessage := func(s *discord.Session, m *discord.MessageCreate){
		if m.Author.ID == s.State.User.ID {
			return
		}
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			if channel, err = s.Channel(m.ChannelID); err != nil {
				eror(err)
			}
		}

		if strings.Contains(m.Content, "!wordx ") {
			x, _ = strconv.Atoi(strings.Replace(m.Content, "!wordx ", "", -1))
		}
		if strings.Contains(m.Content, "!wordy ") {
			y, _ = strconv.Atoi(strings.Replace(m.Content, "!wordy ", "", -1))
		}
		if strings.Contains(m.Content, "!word ") {
			s.ChannelMessageSend(m.ChannelID, "```"+generateBoard(x, y, strings.Replace(m.Content, "!word ", "", -1))+"```")
		}
	
		user := 4
		if m.Author.ID == botuser {
			user += 4
		}

		mentioned := false
		for i := 0; i < len(m.Mentions); i++ {
			if m.Mentions[i].ID == s.State.User.ID {
				mentioned = true
			}
		}
	
		logger.Println(m.Content)
		if channel.Type == discord.ChannelTypeDM && channel.ID == dm.ID {
			fmt.Print("\nDM: "+m.Content+"\n"+">>> ")
		}
		if mentioned || rand.Intn(100) <= user || channel.Type == discord.ChannelTypeDM {
			if message != "null" {
				//s.ChannelMessageSendReply(m.ChannelID, message, (*m).Reference())
			}
		}
	}
	dg.AddHandler(newDiscordMessage)
}

func eror(err error) {
	if err != nil {
		logger.Println(err)
	}
}