package lib

import (
	"time"

	"github.com/briandowns/spinner"
)

func CreateSpinner(prefix, finMessage string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = prefix + " "
	s.FinalMSG = finMessage + "\n"
	s.Start()
	return s
}
