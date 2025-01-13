package utils

import (
	"biFebriansyah/gostream/repositories"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/exp/rand"
)

var birates = []string{"128k", "64k", "32k"}
var progress *Progress = nil
var progresStar = false

const encodeOutFolder string = "./output/encode/"
const segmenOutFolder string = "./output/segment/"

type outStructure struct {
	uid      string
	name     string
	folder   string
	Location string
	bitrate  string
}

func FFmpegStatus() fiber.Handler {
	return websocket.New(func(ctx *websocket.Conn) {
		var currentNum = -1
		for progresStar && !progress.onFinish {
			if currentNum < progress.Value {
				currentNum = progress.Value
				if err := ctx.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(progress.Value))); err != nil {
					log.Println("Failed to send progress:", err)
					break
				}
			}
		}
	})
}

func FFmpegExexute(uid, location string) {
	upload := NewGIO()
	database := NewDatabase()
	repos := repositories.NewMusic(database.DB)
	progress = ProgressInit(100)
	progresStar = true

	encodeAudio := EncodedAudio(location)
	segmentAudio := SegmentedAudio(encodeAudio)
	locations := GenerateMasterPlaylist(segmentAudio)
	if url, err := upload.UploadFolder(locations); err == nil {
		progress.Add(5)
		_, err := repos.InsertSource(url, uid)
		if err != nil {
			log.Println(err)
		}

		progress.Finish()
		log.Printf(`update source url uid: %s`, uid)
	} else {
		log.Printf(`fail to upload with err: %s `, err.Error())
	}
}

func GenerateMasterPlaylist(data *[]outStructure) string {
	var saveLine []string
	var folderUID string
	saveLine = append(saveLine, "#EXTM3U", "#EXT-X-VERSION:3")

	for _, v := range *data {
		readFile, err := os.Open(v.Location)
		if err != nil {
			fmt.Println(err)
		}
		progress.Add(2)

		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		progress.Add(2)

		pattern, err := regexp.Compile("(EXTM3U|EXT-X-VERSION)")
		if err != nil {
			fmt.Println(err)
		}

		for fileScanner.Scan() {
			lines := fileScanner.Text()
			if !pattern.MatchString(lines) && lines != "" {
				saveLine = append(saveLine, fileScanner.Text())
			}
		}
		progress.Add(2)

		masterFile, err := os.Create(v.folder + "master.m3u8")
		if err != nil {
			fmt.Println(err)
		}
		progress.Add(2)

		for _, v := range saveLine {
			_, err := masterFile.WriteString(fmt.Sprint(v + "\n"))
			if err != nil {
				fmt.Println(err)
				masterFile.Close()
			}
		}

		masterFile.Close()
		readFile.Close()
		folderUID = v.folder
		progress.Add(2)
	}

	return folderUID

}

func EncodedAudio(inFileName string) *[]outStructure {
	var outFile []outStructure
	foderUID := uuid.New().String()

	for _, v := range birates {
		uidds := uuid.New().String()
		fileName := uidds[:9] + v + ".aac"
		folderName := encodeOutFolder + foderUID[:8] + "/"
		progress.Add(2)

		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			err := os.Mkdir(folderName, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		progress.Add(2)

		outName := folderName + fileName
		err := ffmpeg.Input(inFileName).
			Output(outName, ffmpeg.KwArgs{"b:a": v, "c:a": "aac"}).
			Silent(true).
			Run()

		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		progress.Add(4)

		outFile = append(outFile, outStructure{
			uid:      foderUID[:8],
			name:     fileName,
			Location: outName,
			bitrate:  v,
		})
		progress.Add(2)
	}

	return &outFile
}

func SegmentedAudio(inStructure *[]outStructure) *[]outStructure {
	var outFile []outStructure
	for _, v := range *inStructure {
		folderName := segmenOutFolder + v.uid + "/"
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			progress.Add(2)
			err := os.Mkdir(folderName, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		progress.Add(2)

		outName := folderName + v.name + ".m3u8"
		masterPl := "master_" + v.bitrate + ".m3u8"
		err := ffmpeg.Input(v.Location).
			Output(outName, ffmpeg.KwArgs{
				"c":                 "copy",
				"f":                 "hls",
				"hls_time":          "10",
				"hls_playlist_type": "vod",
				"master_pl_name":    masterPl,
			}).
			Silent(true).
			Run()

		progress.Add(4)
		if err != nil {
			panic(err)
		}

		outFile = append(outFile, outStructure{
			name:     masterPl,
			folder:   folderName,
			Location: folderName + masterPl,
			bitrate:  v.bitrate,
		})
		progress.Add(2)
	}
	return &outFile
}

func ExampleShowProgress(inFileName, outFileName string) {
	a, err := ffmpeg.Probe(inFileName)
	if err != nil {
		panic(err)
	}

	totalDuration := gjson.Get(a, "format.duration").Float()
	err = ffmpeg.Input(inFileName).
		Output(outFileName, ffmpeg.KwArgs{"b:a": "64k"}).
		GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
		OverWriteOutput().
		Run()

	if err != nil {
		panic(err)
	}
}

func TempSock(totalDuration float64) string {
	rand.Seed(uint64(time.Now().Unix()))
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := 0
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := 0
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = int(float64(c) / (totalDuration * 1000000) * 100)
			}
			if strings.Contains(data, "progress=end") {
				cp = 100
			}

			if cp > progress && cp <= 100 {
				progress = cp
				fmt.Printf("Progress: %d%%\n", progress)
			}

		}
	}()

	return sockFileName
}
