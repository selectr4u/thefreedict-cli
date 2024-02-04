package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	//	"net/http"
)

const (
	reset        string = "\033[0m"
	freedict_url string = "https://api.dictionaryapi.dev/api/v2/entries/en/"
)

// styles
const (
	none int = iota
	bold
	dim
	itallic
	underlined
)

// background colours
const (
	none_bg  int = 0
	black_bg int = iota + 39
	red_bg
	green_bg
	yellow_bg
	blue_bg

	white_bg int = 47
)

// foreground colours
const (
	none_fg  int = 0
	black_fg int = iota + 29
	red_fg
	green_fg
	yellow_fg
	blue_fg

	white_fg int = 37
)

var (
	word string
)

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Phonetics struct {
	Text      string  `json:"text"`
	Audio     string  `json:"audio"`
	SourceURL string  `json:"sourceUrl"`
	License   License `json:"license"`
}

type Definitions struct {
	Definition string   `json:"definition"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
	Example    string   `json:"example"`
}

type Meanings struct {
	PartOfSpeech string        `json:"partOfSpeech"`
	Definitions  []Definitions `json:"definitions"`
	Synonyms     []string      `json:"synonyms"`
	Antonyms     []string      `json:"antonyms"`
}

type DictionaryResponse struct {
	Word       string      `json:"word"`
	Phonetics  []Phonetics `json:"phonetics"`
	Meanings   []Meanings  `json:"meanings"`
	License    License     `json:"license"`
	SourceURLs []string    `json:"sourceUrls"`
}

func init() {
	flag.StringVar(&word, "word", "hello", "the word to search the dictionary for")
}

func sendDictionaryRequest(word string) ([]DictionaryResponse, error) {
	var dictionaryResponses []DictionaryResponse
	response, err := http.Get(fmt.Sprintf("%s%s", freedict_url, word))
	if err != nil {
		fmt.Println("Error sending GET request to API:", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", response.StatusCode)
		return nil, fmt.Errorf("request failed with status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	err = json.Unmarshal(body, &dictionaryResponses)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	return dictionaryResponses, nil
}

// specific formatting functions

func textFormat(str string, styles ...int) string {

	var stylesString string = "\033["

	for i := 0; i < len(styles); i++ {
		stylesString = fmt.Sprintf("%s;%d", stylesString, styles[i])

		if i == len(styles)-1 {
			stylesString = fmt.Sprintf("%sm", stylesString)
		}

	}

	return fmt.Sprintf("%s%s%s", stylesString, str, reset)
}

func formatMeanings(meanings []Meanings) (string, error) {
	var stringOutput string = ""

	for _, meaning := range meanings {
		posString := fmt.Sprintf("Part of Speech: %s", meaning.PartOfSpeech)
		definitionString, err := formatDefinitions(meaning.Definitions)

		if err != nil {
			return "", err
		}

		stringOutput = fmt.Sprintf("%s\n%s\n\n", posString, definitionString)

	}

	return stringOutput, nil
}

func formatDefinitions(definitions []Definitions) (string, error) {
	var stringOutput string = ""

	for i := 0; i < len(definitions); i++ {
		var definitionString string
		var definitionTextBoldTitle string = textFormat(fmt.Sprintf("Definition %d:", i+1), underlined)
		var exampleTextTitle string = textFormat("Example:", underlined)

		definitionString = fmt.Sprintf("  %s \n  %s", definitionTextBoldTitle, definitions[i].Definition)
		if definitions[i].Example != "" {
			definitionString = fmt.Sprintf("%s \n  %s", definitionString, fmt.Sprintf("%s %s", exampleTextTitle, definitions[i].Example))
		}

		definitionString += "\n"

		stringOutput = fmt.Sprintf("%s%s", stringOutput, definitionString)

	}

	return stringOutput, nil
}

func formatDictionaryResponse(response *DictionaryResponse) (string, error) {
	var formattedString, formattedMeanings string

	formattedString = textFormat(fmt.Sprintf("%s (phonetic: '%s')", response.Word, response.Phonetics[0].Text), bold, underlined, itallic)

	formattedString = fmt.Sprintf("%s\n\n%s\n", formattedString, textFormat("Meanings:", bold, underlined))

	formattedMeanings, err := formatMeanings(response.Meanings)

	if err != nil {
		return "", err
	}

	formattedString = fmt.Sprintf("%s%s\n", formattedString, formattedMeanings)

	return formattedString, nil

}

func main() {
	flag.Parse()

	fmt.Println("Searching dictionary for", word)

	dictionaryResponses, err := sendDictionaryRequest(word)
	if err != nil {
		fmt.Println("Unable to complete request due to an error:", err)
		return
	}

	for _, dictionaryResponse := range dictionaryResponses {
		dictionaryString, err := formatDictionaryResponse(&dictionaryResponse)
		if err != nil {
			fmt.Println("Unable to format dictionary response:", err)
			continue
		}
		fmt.Println(dictionaryString)
	}
}
