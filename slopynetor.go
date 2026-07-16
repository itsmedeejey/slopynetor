package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func main() {
	var url string
	fmt.Print("Enter the youtube URL: ")

	if _, err := fmt.Scanf("%s", &url); err != nil {
		fmt.Println("error reading input:", err)
		return
	}

	// download(url)
	// toAudio()
	transcripe()
}
