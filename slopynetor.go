package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func download(url string) {
	// creats the download folder if doesn't exist
	if err := os.MkdirAll("download", 0755); err != nil {
		fmt.Printf("folder creation failed: %v\n", err)
	}

	output := filepath.Join(
		"download",
		"video.%(ext)s",
	)

	cmd := exec.Command("yt-dlp",
		"-o",
		output,
		url,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Download failed: %v\n", err)
		return
	}
}

func toAudio() {
	// here we are reading the download folder and geting the first file
	files, err := os.ReadDir("download")
	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		fmt.Println("No files found")
		return
	}
	file := files[0]

	videoPath := filepath.Join(
		"download",
		file.Name(),
	)

	audioPath := filepath.Join(
		"download",
		"audio.mp3",
	)

	cmd := exec.Command(
		"ffmpeg",
		"-i",
		videoPath,
		"-vn",
		audioPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("audio conversation failed: %v\n", err)
		return
	}

}

func transcripe() {
	if err := os.MkdirAll("output", 0755); err != nil {
		fmt.Printf("folder creation failed: %v\n", err)
	}

	cmd := exec.Command(
		"trans/.venv/bin/python",
		"trans/transcribe.py",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("transcription failed: %v\n", err)
		return
	}

}

func getClips() {
	// calling gemini for the trendy clips from the transcription
	content, err := os.ReadFile("output/transcript.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	promptContent, err := os.ReadFile("podcastPrompt.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert bytes to string and print
	transcript := string(content)
	promptText := string(promptContent)

	prompt := fmt.Sprintf(`%s transcript: %s`, promptText, transcript)

	client, err := genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-3.5-flash", //make it env variable
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Text())
	if err := os.WriteFile("output/clips.json", []byte(resp.Text()), 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File saved successfully!")
}

type Clip struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Title string  `json:"title"`
}

func cuts() {

	videoPath := filepath.Join(
		"download",
		"video.webm",
	)

	//geting  the json file
	data, err := os.ReadFile("output/clips.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	//parsing the json file
	var clips []Clip
	if err = json.Unmarshal(data, &clips); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	exportDict := fmt.Sprintln("output/shorts")
	if err := os.MkdirAll(exportDict, 0755); err != nil {
		fmt.Printf("folder creation failed: %v\n", err)
	}

	for _, clip := range clips {

		startStr := strconv.Itoa(int(clip.Start))
		endStr := strconv.Itoa(int(clip.End))

		outputFile := fmt.Sprintf("%s/%s.mp4", exportDict, clip.Title)
		cmd := exec.Command(
			"ffmpeg", "-y",
			"-ss", startStr,
			"-to", endStr,
			"-i", videoPath,
			"-vf", "crop=ih*4/3:ih,scale=1080:-2,pad=1080:1920:(ow-iw)/2:(oh-ih)/2",
			"-c:v", "libx264",
			"-preset", "medium",
			"-crf", "23",
			"-c:a", "aac",
			"-b:a", "128k",
			outputFile,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("transcription failed: %v\n", err)
			return
		}
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var url string
	fmt.Print("Enter the youtube URL: ")

	if _, err := fmt.Scanf("%s", &url); err != nil {
		fmt.Println("error reading input:", err)
		return
	}

	download(url)
	toAudio()
	transcripe()
	getClips()
	cuts()
	fmt.Println("completed!")
}
