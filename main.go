package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"BadBot/telegram"
	"BadBot/lib"

	"io/ioutil"
	//"regexp"
)



func init() {
	flag.StringVar(&lib.TelegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()
	if lib.TelegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
}

func main() {
	logs := make(chan string)
	go telegram.MainTtelegram(logs)
	handleHello := makeHello(logs)
	http.HandleFunc("/gateway/telegram/create/bad", handleHello)
	http.ListenAndServe(":8181", nil)

}

func makeHello(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "all ok	")
		decoder := json.NewDecoder(r.Body)
		var t lib.BadJson

		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(ioutil.ReadAll(r.Body))
			fmt.Fprint(w, "all ok	")
			log.Println(err)
			return
		}
		string := fmt.Sprint("id: ", t.Ip, "info: ", t.Json)
				fmt.Println("Отправил в логгер")
		select {
		case  <-logger:
			logger <- string
		default:
			return
		}
	}
}