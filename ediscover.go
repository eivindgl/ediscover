package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dixonwille/skywalker"
	"github.com/urfave/cli"
)

type ExampleWorker struct {
	*sync.Mutex
	found []string
}

func (ew *ExampleWorker) Work(path string) {
	//This is where the necessary work should be done.
	//This will get concurrently so make sure it is thread safe if you need info across threads.
	ew.Lock()
	defer ew.Unlock()
	if "help.go" == filepath.Base(path) {
		ew.found = append(ew.found, path)
	}
}

func ExampleSkywalker(path string) {
	//Following two functions are only to create and destroy data for the example
	// defer teardownData()
	// standupData()

	ew := new(ExampleWorker)
	ew.Mutex = new(sync.Mutex)

	//root is the root directory of the data that was stood up above
	sw := skywalker.New(path, ew)
	sw.DirListType = skywalker.LTBlacklist
	sw.DirList = []string{"sub"}
	sw.ExtListType = skywalker.LTWhitelist
	sw.ExtList = []string{".go"}
	err := sw.Walk()
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Sort(sort.StringSlice(ew.found))
	for _, f := range ew.found {
		show := strings.Replace(f, sw.Root, "", 1)
		show = strings.Replace(show, "\\", "/", -1)
		fmt.Println(show)
	}
	// Output:
	// /subfolder/few.pdf
	// /the/few.pdf
}

func main() {
	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "fight the loneliness!"
	app.Version = "18.04"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q", c.Args().Get(0))
		fmt.Println("Running skywalker")
		ExampleSkywalker("/Users/gardl/go/src")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
