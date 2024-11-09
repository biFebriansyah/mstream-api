package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	GIOS_BUCKET = "streamapp"
)

type GioStores struct {
	*s3.Client
}

type uploadFD struct {
	name     string
	location string
}

func NewGIO() *GioStores {
	GIOS_URL := os.Getenv("GIOS_URL")
	GIOS_KEY := os.Getenv("GIOS_KEY")
	GIOS_SECRET := os.Getenv("GIOS_SECRET")
	GIOS_REGION := os.Getenv("GIOS_REGION")

	client := s3.New(s3.Options{
		Region:       GIOS_REGION,
		BaseEndpoint: aws.String(GIOS_URL),
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(GIOS_KEY, GIOS_SECRET, "")),
	})

	return &GioStores{client}
}

func GetPublicUrl(key string) string {
	url := fmt.Sprintf("https://%s.nos.jkt-1.neo.id/%s", GIOS_BUCKET, key)
	return url
}

func GetFolderPath(filePath string) string {
	var folderName string
	dirPath := filepath.Dir(filePath)
	pathParts := strings.Split(dirPath, string(filepath.Separator))

	if len(pathParts) >= 3 {
		folderName = pathParts[len(pathParts)-1]
	} else {
		fmt.Println("Unexpected path structure")
	}

	return folderName
}

func (c *GioStores) UploadData(locations string) (string, error) {
	fileName := filepath.Base(locations)

	file, err := os.Open(locations)
	if err != nil {
		return "", err
	}

	_, err = c.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(GIOS_BUCKET),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.nos.jkt-1.neo.id/%s", GIOS_BUCKET, fileName)

	return url, file.Close()
}

func (c *GioStores) ListData() ([]types.Object, error) {
	output, err := c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(GIOS_BUCKET),
	})
	if err != nil {
		return nil, err
	}

	return output.Contents, nil
}

func (c *GioStores) UploadFolder(location string) (string, error) {
	var folderPath string
	err := filepath.Walk(location, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		folderPath = GetFolderPath(path)
		relativePath := folderPath + "/" + fileName

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = c.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(GIOS_BUCKET),
			Key:    aws.String(relativePath),
			Body:   file,
			ACL:    types.ObjectCannedACLPublicRead,
		})

		if err != nil {
			return err
		}

		// fmt.Printf("Uploaded %s to %s\n", path, relativePath)

		return nil

	})

	url := fmt.Sprintf("https://%s.nos.jkt-1.neo.id/%s/%s", GIOS_BUCKET, folderPath, "master.m3u8")
	return url, err
}
