package main

import (
	"biFebriansyah/gostream/routers"
	"biFebriansyah/gostream/utils"
	"context"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	database := utils.NewDatabase()
	routers := routers.New(database.DB)
	server := utils.NewServer(routers)

	wait := utils.GracefulShutdown(context.Background(), 2*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return database.Close()
		},
		"server": func(ctx context.Context) error {
			return server.Shutdown()
		},
	})

	<-wait
}

func tesst() {
	database := utils.NewDatabase()
	routers := routers.New(database.DB)
	server := utils.NewServer(routers)

	wait := utils.GracefulShutdown(context.Background(), 2*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return database.Close()
		},
		"server": func(ctx context.Context) error {
			return server.Shutdown()
		},
	})

	<-wait
}

// func generateAudioStream() {
// 	encodeAudio := EncodedAudio("input_test.mp3")
// 	segmentAudio := SegmentedAudio(encodeAudio)
// 	location := GenerateMasterPlaylist(segmentAudio)

// 	gioStore := NewGIO()
// 	data, err := gioStore.UploadFolder(location)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(data)
// }
