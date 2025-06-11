package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	InputFile   string            `json:"InputFile"`
	OutputFile  string            `json:"OutputFile"`
	Title       string            `json:"Title"`
	Author      string            `json:"Author"`
	FontPresets map[string]string `json:"FontPresets"`
}

func writeExampleFile() {
	// example
	exampleConfigContent := []byte(`{
	"InputFile":"files/hw3.tex",
	"OutputFile": "files/hw3-hants.tex",
	"Title": "人工智慧概論 HW3",
	"Author": "4112064214 侯竣奇",
	"FontPresets": {
		"setmainfont": "Times New Roman",
		"setmonofont": "CaskaydiaCove Nerd Font Mono",
		"setCJKmainfont": "標楷體"
	}
}`)

	// 寫入 config.json 檔
	err := os.WriteFile("config.json", []byte(exampleConfigContent), 0644)
	if err != nil {
		log.Fatalf("Error writing config.json: %v", err)
	}
}

func main() {
	// 讀取 config.json 檔
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		// 如果找不到 config.json 檔，就建立一個範例檔案
		if os.IsNotExist(err) {
			writeExampleFile()
			log.Println("config.json not found, created a example file")
			return
		}
		// 其他錯誤就直接報錯
		log.Fatalf("Error reading config.json, please check the file: %v", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error parsing config.json, please check the file: %v. You can delete the config.json and run again to create a example file.", err)
	}

	// 依照 InputFile 讀取 .tex 檔
	inputFile, err := os.ReadFile(config.InputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// 依照 config 裡的設定，建立 preamble
	preambles := []string{
		"    \\usepackage{fontspec}",
		"    \\usepackage[slantfont, boldfont]{xeCJK}",
	}
	for key, value := range config.FontPresets {
		preambles = append(preambles, fmt.Sprintf("    \\%s{%s}", key, value))
	}

	// 在 .tex 檔中第二行的位置插入 preamble
	inputFileLines := strings.Split(string(inputFile), "\n")
	outputFileLines := []string{inputFileLines[0]}
	outputFileLines = append(outputFileLines, preambles...)
	outputFileLines = append(outputFileLines, inputFileLines[1:]...)

	var titleLine int
	for i, line := range outputFileLines {
		if strings.Contains(line, "    \\title{") {
			titleLine = i
			break
		}
	}
	outputFileLines[titleLine] = fmt.Sprintf("    \\title{%s}", config.Title)

	// 在 title 後面插入 author
	beforeAuthor := make([]string, titleLine+1)
	copy(beforeAuthor, outputFileLines[:titleLine+1])
	afterAuthor := make([]string, len(outputFileLines)-titleLine-1)
	copy(afterAuthor, outputFileLines[titleLine+1:])
	outputFileLines = append(beforeAuthor, fmt.Sprintf("    \\author{%s}", config.Author))
	outputFileLines = append(outputFileLines, afterAuthor...)

	// 寫入 outputFileLines
	err = os.WriteFile(config.OutputFile, []byte(strings.Join(outputFileLines, "\n")), 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

}
