package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apiKey = "YOUR API_KEY"
const searchURL = "https://translation.googleapis.com/language/translate/v2"
var targetLang string

type TransResp struct {
	Data struct {
		Translations []Translation	`json:"Translations"`
	}
}

type Translation struct {
	TranslatedText string	`json:"translatedText"`
	DetectedLang string	`json:"detectedSourceLanguage"`
}

func getValues(args []string) url.Values {
	values := url.Values{}
	values.Add("q", strings.Join(args[1:], " "))
	values.Add("target", args[0])
	values.Add("format", "text")
	values.Add("key", apiKey)
	return values
}

func execute(res *http.Response) {
	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ts TransResp
	err2 := json.Unmarshal(byteArray, &ts)
	if err2 != nil {
		fmt.Println(err2)
	}

	display(ts)
}

func display(ts TransResp) {
	for _, t := range ts.Data.Translations {
		fmt.Printf("%s->%s: %s", t.DetectedLang, targetLang, t.TranslatedText)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage: Enter the target language code and the search words
        [Example] ja wehre are you from?`)
	}
	flag.Parse()

	targetLang = flag.Args()[0]
	values := getValues(flag.Args())
	res, err := http.Get(searchURL + "?" + values.Encode())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	execute(res)
}
