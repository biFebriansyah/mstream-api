package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CloudInary(file interface{}) (string, error) {
	name := os.Getenv("CD_NAME")
	key := os.Getenv("CD_KEY")
	secret := os.Getenv("CD_SECRET")

	cld, err := cloudinary.NewFromParams(name, key, secret)
	if err != nil {
		return "", err
	}

	uploadParams := &uploader.UploadParams{
		ResourceType: "video",
		Folder:       "music/hls/test",
		Eager:        "ac_aac,af_96000,br_3500k,q_auto,vc_h264/m3u8",
		EagerAsync:   api.Bool(true),
	}

	result, err := cld.Upload.Upload(context.Background(), file, *uploadParams)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return result.URL, nil
}
