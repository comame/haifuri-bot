package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
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

var diags = []string{
	"海に生き",
	"海を守り",
	"海を往く",
	"それがブルーマーメイド！",
	"海の仲間は、家族だから",
	"艦長、私はあなたのマヨネーズになる！",
	":muscle:ハイスクール・コマンドー:muscle:",
	":tada:ハイスクール・フリート:tada:",
	"ん、私ともあっちむいてホイをやるべき",
	"行こう行こう、いつも先を急ぐ。⊂二二⌒°( ＾ω＾)°⌒二二⊃ そしてある日死ぬ。",
}

func shuffledHaifuri() string {
	percent := rand.Int31n(100)
	if percent < 10 {
		i := rand.Int31n(int32(len(diags)))
		return diags[i]
	}

	original := "ハイスクール・フリート"
	chars := []rune(original)

	rand.Shuffle(len(chars), func(i, j int) {
		tmp := chars[i]
		chars[i] = chars[j]
		chars[j] = tmp
	})

	first := chars[0]
	end := chars[len(chars)-1]

	for first == '・' || first == 'ー' || end == '・' {
		rand.Shuffle(len(chars), func(i, j int) {
			tmp := chars[i]
			chars[i] = chars[j]
			chars[j] = tmp
		})

		first = chars[0]
		end = chars[len(chars)-1]
	}

	if string(chars) == "ハイスクール・フリート" {
		return ":tada:" + string(chars) + ":tada:"
	}

	return string(chars)
}

func main() {
	checkEnv()

	rand.Seed(time.Now().UnixNano())
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
					"content": shuffledHaifuri(),
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
		fmt.Fprintf(w, string(resStr))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("not found")
		http.Error(w, "not found: "+r.URL.Path, 404)
	})
	println("Start")
	http.ListenAndServe(":8080", nil)
}
