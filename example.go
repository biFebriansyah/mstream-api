package main

import "fmt"

type arrObj map[string]map[string]string

type test struct {
	name   string
	hoby   []string
	schooL arrObj
	email  string
}

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width, Height float64
}

type Triangle struct {
	Base, Height float64
}

func example() {
	var p test
	p.name = "ebiebi"
	p.hoby = []string{"asas", "asas"}
	p.schooL = arrObj{
		"smk": {
			"name": "",
			"majo": "sas",
		},
	}

	fmt.Println(p)
}
