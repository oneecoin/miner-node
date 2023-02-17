package properties

import "github.com/fatih/color"

var ErrorStr func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()
var WarningStr func(a ...interface{}) string = color.New(color.FgYellow).SprintFunc()
