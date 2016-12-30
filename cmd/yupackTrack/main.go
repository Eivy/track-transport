package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"local/transportTrack/yupack"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	var interactive bool
	var pipe bool
	var csv bool
	flag.BoolVar(&interactive, "i", false, "use interactive")
	flag.BoolVar(&pipe, "p", false, "use pipe")
	flag.BoolVar(&csv, "csv", false, "output csv(Shift-JIS, CRLF)")
	flag.Parse()
	if interactive && pipe {
		fmt.Fprintln(os.Stderr, "Can not use interactive mode and pipe mode same time")
		os.Exit(99)
	}
	if interactive {
		fmt.Print("Input tracking number: ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if v, err := yupack.GetStatus([]string{scanner.Text()}); err != nil {
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
		if v, err := yupack.GetStatus(strings.Split(string(t), "\n")); err != nil {
			exitWithError(err)
		} else {
			output(v, csv)
		}
	} else {
		for _, v := range flag.Args() {
			b, err := ioutil.ReadFile(v)
			if err != nil {
				exitWithError(err)
			}
			if v, err := yupack.GetStatus(strings.Split(string(b), "\n")); err != nil {
				exitWithError(err)
			} else {
				output(v, csv)
			}
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func output(value [][]string, csvFlag bool) {
	if csvFlag {
		file, err := os.Create("result.csv")
		if err != nil {
			exitWithError(err)
		}
		defer file.Close()
		w := csv.NewWriter(transform.NewWriter(file, japanese.ShiftJIS.NewEncoder()))
		w.UseCRLF = true
		for _, line := range value {
			w.Write(line)
		}
		w.Flush()
	} else {
		for _, n := range value {
			fmt.Println(strings.Join(n, ","))
		}
	}
}
