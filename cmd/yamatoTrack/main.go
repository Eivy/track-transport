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
		fmt.Print("Input tracking number: ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if v, err := yamato.GetStatus([]string{scanner.Text()}); err != nil {
				exitWithError(err)
			} else {
				for _, value := range v {
					fmt.Println(strings.Join(value, ","))
				}
			}
			fmt.Print("Input tracking number: ")
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
				if v, err := yamato.GetStatus(strings.Split(string(b), "\n")); err != nil {
					exitWithError(err)
				} else {
					for _, n := range v {
						fmt.Println(strings.Join(n, ","))
					}
				}
			}
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
