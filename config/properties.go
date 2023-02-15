package config

import "github.com/fatih/color"

var PublicKey string
var Port int

var ErrorStr func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()
