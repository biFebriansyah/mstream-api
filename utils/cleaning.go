package utils

import (
	"fmt"
	"log"
	"os"
)

type cleaning struct {
	Data []string
}

func Cleaning() *cleaning {
	return &cleaning{}
}

func (c *cleaning) Add(data ...string) {
	for _, v := range data {
		c.Data = append(c.Data, v)
	}
}

func (c *cleaning) Run() {
	for _, v := range c.Data {
		fileInfo, err := os.Stat(v)
		if err != nil {
			log.Println(err)
		}

		if fileInfo.IsDir() {
			if err := os.RemoveAll(v); err != nil {
				fmt.Printf("Failed to remove file %s: %v\n", v, err)
			}
		} else {
			if err := os.Remove(v); err != nil {
				fmt.Printf("Failed to remove file %s: %v\n", v, err)
			}
		}

	}
}
