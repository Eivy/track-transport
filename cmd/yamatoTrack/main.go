package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"local/transportTrack/yamato"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if err := yamato.GetStatus([]string{scanner.Text()}); err != nil {
				exitWithError(err)
			}
		}
		if err := scanner.Err(); err != nil {
			exitWithError(err)
		}
	} else {
		for i, v := range os.Args {
			if i > 0 {
				b, err := ioutil.ReadFile(v)
				if err != nil {
					exitWithError(err)
				}
				yamato.GetStatus(strings.Split(string(b), "\n"))
			}
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
