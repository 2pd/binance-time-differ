package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

type Flag struct {
	verbose bool
}

func getLocalTime() int {
	now := time.Now()
	ms := now.UnixNano() / 1000000
	return int(ms)
}

func getServerTime(f *Flag) int {

	url := "https://api.binance.com/api/v3/time"

	if f.verbose {
		fmt.Printf("Send request to : %#v\n", url)
	}
	resp, _ := http.Get(url)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	data := string(body)
	time, err := strconv.Atoi(data[14:27])

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if f.verbose {
		fmt.Printf("Get server time: %#v\n", time)
	}
	return time

}

func timeDiff(t1, t2 int) int {
	return t1 - t2
}

func localvsServer(f *Flag) {

	st := getServerTime(f)
	t2 := getLocalTime()
	if f.verbose {
		fmt.Printf("Local time: %#v\n", t2)
	}
	fmt.Println(timeDiff(st, t2))
}

func main() {
	app := &cli.App{
		Name:  "Binance Time Differ",
		Usage: "Calculate the time diff from Binance Server to locally",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
		},
		Action: func(c *cli.Context) error {
			var f Flag

			f.verbose = c.Bool("v")

			localvsServer(&f)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
