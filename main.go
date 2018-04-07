package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
)

type BotInfo struct {
	ClientId	string `json:"client_id"`
	Token 		string `json:"token"`
}

var (
	stopBot    = make(chan bool)
	HelloWorld = "!hello"
)

func initialize(botInfo BotInfo) (botName string, botToken string) {
	botName = "<@"+string(botInfo.ClientId)+">"
	botToken = "Bot "+string(botInfo.Token)
	return
}

func main() {
	//BotInfoにクライアントIDとトークンを突っ込む
	var botInfo BotInfo
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err2 := json.Unmarshal(file, &botInfo)
	if err2 != nil {
		log.Fatal(err2)
	}
	_, Token := initialize(botInfo)

	//Discordのセッションを作成
	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加
	// websocketを開いてlistening開始
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	<-stopBot //プログラムが終了しないようロック
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//BotInfoにクライアントIDとトークンを突っ込む
	var botInfo BotInfo
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err2 := json.Unmarshal(file, &botInfo)
	if err2 != nil {
		log.Fatal(err2)
	}
	BotName, _ := initialize(botInfo)
	c, err := s.State.Channel(m.ChannelID) //チャンネル取得
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", BotName, HelloWorld)):
		sendMessage(s, c, "私は桜ねね！！\nよろしく！！")
	}
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
