package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var APPLICATION_ID = os.Getenv("APPLICATION_ID")
var GUILD_ID = os.Getenv("GUILD_ID")
var PUBLIC_KEY = os.Getenv("PUBLIC_KEY")
var BOT_TOKEN = os.Getenv("BOT_TOKEN")

var COMMANDS_URL = fmt.Sprintf("https://discord.com/api/v10/applications/%s/commands", APPLICATION_ID)

func checkEnv() {
	if APPLICATION_ID == "" || GUILD_ID == "" || PUBLIC_KEY == "" || BOT_TOKEN == "" {
		panic("env not present")
	}
}

func makeCommandApi(command *Commands) {
	body, err := json.Marshal(command)
	if err != nil {
		panic(err)
	}

	bodyReader := strings.NewReader(string(body))

	req, _ := http.NewRequest("POST", COMMANDS_URL, bodyReader)
	req.Header.Set("authorization", fmt.Sprintf("Bot %s", BOT_TOKEN))
	req.Header.Set("content-type", "application/json")

	res, err := new(http.Client).Do(req)
	if err != nil {
		panic(err)
	}

	resStr, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resStr))
}

func makeCommand() {
	command := Commands{
		Name:        "haifuri",
		Type:        1,
		Description: "ハイスクール・フリート",
	}
	makeCommandApi(&command)
}

// [min, max] の範囲でランダムな整数を返す
func randomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func shuffleHaifuri() string {
	// カタカナ8文字をシャッフル
	var shuffledKatakana = []rune("ハイスクルフリト")
	rand.Shuffle(len(shuffledKatakana), func(i, j int) {
		shuffledKatakana[i], shuffledKatakana[j] = shuffledKatakana[j], shuffledKatakana[i]
	})

	// 丸の位置を決める (カタカナの間に入る)
	var dotPos = randomInt(1, 6)

	// 伸ばし棒をつけるカタカナを8文字から2つ選ぶ
	var firstHyphenPos int
	var secondHyphenPos int
	for firstHyphenPos == secondHyphenPos {
		firstHyphenPos = randomInt(0, 7)
		secondHyphenPos = randomInt(0, 7)
	}

	var result string
	// 組み立て
	for i := 0; i <= 7; i++ {
		if i == dotPos {
			result += "・"
		}
		result += string(shuffledKatakana[i])
		if firstHyphenPos == i || secondHyphenPos == i {
			result += "ー"
		}
	}

	// 一致したらお祝い
	if result == "ハイスクール・フリート" {
		return ":tada:" + result + ":tada:"
	}

	return result
}

func main() {
	log.Println("Built 222643")

	checkEnv()

	makeCommand()

	http.HandleFunc("/high-school-fleet/interactions", func(w http.ResponseWriter, r *http.Request) {
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		sigEd25519 := r.Header.Get("X-Signature-Ed25519")
		sigTimestamp := r.Header.Get("X-Signature-Timestamp")

		verified := verifySignature(sigEd25519, sigTimestamp, string(rawBody))
		if !verified {
			println("unauthorized")
			http.Error(w, "unauthorized", 403)
			return
		}

		var req InteractionRequest
		err = json.Unmarshal(rawBody, &req)
		if err != nil {
			http.Error(w, "error", 400)
			return
		}

		var res InteractionResponse
		switch req.Type {
		case 1:
			res = InteractionResponse{
				Type: 1,
				Data: nil,
			}
		case 2:
			res = InteractionResponse{
				Type: 4,
				Data: map[string]interface{}{
					"content": shuffleHaifuri(),
				},
			}
		default:
			println("unrecognized interaction type")
			http.Error(w, "error", 400)
			return
		}

		resStr, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "error", 400)
			return
		}

		println(string(resStr))

		w.Header().Add("content-type", "application/json")
		fmt.Fprintln(w, string(resStr))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("not found")
		http.Error(w, "not found: "+r.URL.Path, 404)
	})
	println("Start")
	http.ListenAndServe(":8080", nil)
}
