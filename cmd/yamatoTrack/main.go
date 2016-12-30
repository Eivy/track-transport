package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"local/transportTrack/yamato"
	"os"
	"strings"
)

func main() {
	var interactive bool
	var pipe bool
	flag.BoolVar(&interactive, "i", false, "use interactive")
	flag.BoolVar(&pipe, "p", false, "use pipe")
	flag.Parse()
	if interactive && pipe {
		fmt.Fprintln(os.Stderr, "Can not use interactive mode and pipe mode same time")
		os.Exit(99)
	}
	if interactive {
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
	} else if pipe {
		t, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			exitWithError(err)
		}
		if v, err := yamato.GetStatus(strings.Split(string(t), "\n")); err != nil {
			exitWithError(err)
		} else {
			for _, n := range v {
				fmt.Println(strings.Join(n, ","))
			}
		}
	} else {
		for _, v := range flag.Args() {
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

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
