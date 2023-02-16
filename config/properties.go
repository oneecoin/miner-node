package config

import "github.com/fatih/color"

const (
	DefaultPageSize = 10
	MempoolAddress  = "localhost:8080"
)

var PublicKey string
var Port int

var ErrorStr func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()
