package main

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	encodeAudio := EncodedAudio("input_test.mp3")
	segmentAudio := SegmentedAudio(encodeAudio)
	location := GenerateMasterPlaylist(segmentAudio)

	gioStore := NewGIO()
	data, err := gioStore.UploadFolder(location)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}

func generateAudioStream() {
	encodeAudio := EncodedAudio("input_test.mp3")
	segmentAudio := SegmentedAudio(encodeAudio)
	location := GenerateMasterPlaylist(segmentAudio)
	fmt.Println(location)
}
