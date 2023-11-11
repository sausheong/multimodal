package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var cyan = color.New(color.FgCyan).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	files := []string{"award.jpg", "fruits.jpg", "vortex.jpg", "fisherman.jpg", "untidy_corridor.jpg"}
	question := "What is in this picture and where is it?"
	for _, file := range files {
		base64Img, err := getBase64("./test-images/" + file)
		if err != nil {
			panic(err)
		}
		fmt.Println()
		// show question and picture
		fmt.Println(cyan(question + " - " + file))
		fmt.Printf("\n \033]1337;File=width=32;inline=1:%s\a\n", base64Img)

		// OpenAI
		fmt.Println(yellow("OpenAI GPT4-Vision"))
		start := time.Now()
		gptvRes, err := gptv("", base64Img)
		if err != nil {
			panic(err)
		}
		fmt.Println(gptvRes)
		fmt.Println(green(time.Since(start).Round(time.Second)))

		// Google Vertex AI
		fmt.Println(yellow("\nGoogle VertexAI Imagen"))
		start = time.Now()
		imagenRes, err := imagen(question, base64Img)
		if err != nil {
			panic(err)
		}
		fmt.Println(imagenRes[0])
		fmt.Println(green(time.Since(start).Round(time.Second)))

		// Llava-1.5-7B
		fmt.Println(yellow("\nLlava-1.5-7B"))
		start = time.Now()
		llavaRes, err := llava(question, "./test-images/"+file)
		if err != nil {
			panic(err)
		}
		fmt.Println(llavaRes)
		fmt.Println(green(time.Since(start).Round(time.Second)))
	}
}

func getBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
