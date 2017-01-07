package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"local/transportTrack/fukutu"
	"local/transportTrack/sagawa"
	"local/transportTrack/yamato"
	"local/transportTrack/yupack"
	"os"
	"strings"

	"github.com/mattn/go-colorable"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var stdOut io.Writer

func main() {
	var pipe, usage bool
	flag.BoolVar(&pipe, "p", false, "use pipe")
	flag.BoolVar(&usage, "h", false, "show this usage")
	flag.Parse()
	stdOut = colorable.NewColorableStdout()
	if usage {
		flag.Usage()
		return
	}
	if pipe && len(flag.Args()) > 0 {
		fmt.Println("Can not read file with pipe flag")
		os.Exit(99)
	}
	if pipe {
		if v, err := read(os.Stdin); err != nil {
			exitWithError(err)
		} else {
			output(v, false)
		}
	} else if len(flag.Args()) == 1 {
		f, err := os.Open(flag.Args()[0])
		if err != nil {
			exitWithError(err)
		}
		if v, err := read(f); err != nil {
			exitWithError(err)
		} else {
			output(v, true)
		}
	} else {
		const coopMsg = "\x1b[36m1:ヤマト運輸; 2:ゆうパック; 3:佐川急便; 4:福山通運; 0: 終了する\x1b[0m"
		const chooseCoopMsg = "\x1b[36m会社を選んでください: \x1b[0m"
		const numberMsg = "\x1b[36m送り状番号: \x1b[0m"
		s := bufio.NewScanner(os.Stdin)
		fmt.Fprintln(stdOut, coopMsg)
		fmt.Fprint(stdOut, chooseCoopMsg)
		for s.Scan() {
			c := s.Text()
			if c == "0" {
				return
			}
			fmt.Fprint(stdOut, numberMsg)
			s.Scan()
			n := s.Text()
			getResultOne(c, n)
			fmt.Fprintln(stdOut, coopMsg)
			fmt.Fprint(stdOut, chooseCoopMsg)
		}
	}
}

func read(reader *os.File) ([][]string, error) {
	r := csv.NewReader(transform.NewReader(reader, japanese.ShiftJIS.NewDecoder()))
	all, err := r.ReadAll()
	if err != nil {
		exitWithError(err)
	}
	c := ""
	ret := [][]string{}
	tmp := []string{}
	for _, record := range all {
		if c != "" && c != record[0] {
			if t, e := getResult(c, tmp); e != nil {
				return nil, e
			} else {
				ret = append(ret, t...)
			}
			tmp = []string{}
		}
		tmp = append(tmp, record[1])
		c = record[0]
	}
	if t, e := getResult(c, tmp); e != nil {
		return nil, e
	} else {
		ret = append(ret, t...)
	}
	return ret, nil
}

func getResultOne(c string, n string) {
	switch c {
	case "1":
		if r, e := yamato.GetStatus([]string{n}); e != nil {
			exitWithError(e)
		} else {
			for _, v := range r {
				fmt.Println(strings.Join(v, ","))
			}
		}
	case "2":
		if r, e := yupack.GetStatus([]string{n}); e != nil {
			exitWithError(e)
		} else {
			for _, v := range r {
				fmt.Println(strings.Join(v, ","))
			}
		}
	case "3":
		if r, e := sagawa.GetStatus([]string{n}); e != nil {
			exitWithError(e)
		} else {
			for _, v := range r {
				fmt.Println(strings.Join(v, ","))
			}
		}
	case "4":
		if r, e := fukutu.GetStatus([]string{n}); e != nil {
			exitWithError(e)
		} else {
			for _, v := range r {
				fmt.Println(strings.Join(v, ","))
			}
		}
	}
}

func getResult(c string, n []string) ([][]string, error) {
	switch c {
	case "1":
		return yamato.GetStatus(n)
	case "2":
		return yupack.GetStatus(n)
	case "3":
		return sagawa.GetStatus(n)
	case "4":
		return fukutu.GetStatus(n)
	}
	return nil, nil
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
