package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/shoukoo/bkb/list"
	"github.com/shoukoo/bkb/screen"
)

var (
	helpBool    bool
	versionBool bool
	version     string
	desc        string
)

func init() {
	version = "0.1.8"
	desc = `
Usage of bbk:
  bbk [flags] # run buildkite beaver

`
	flag.BoolVar(&helpBool, "help", false, "Print help and exist")
	flag.BoolVar(&versionBool, "version", false, "Print version and exit")
}

func run() error {

	stdin := readline.NewCancelableStdin(os.Stdin)
	c := &readline.Config{}
	err := c.Init()
	if err != nil {
		return err
	}

	c.Stdin = stdin

	c.HistoryLimit = -1
	c.UniqueEditLine = true

	rl, err := readline.NewEx(c)
	if err != nil {
		return err
	}

	_, err = rl.Write([]byte(screen.HideCursor))
	if err != nil {
		return err
	}
	t := screen.New(rl)

	client, err := list.BuildkiteClient()
	if err != nil {
		return err
	}
	err = client.GetRecentBuilds("lexer")
	if err != nil {
		return err
	}

	listConfig := list.Config{
		Items:     client.Builds,
		Templates: client.Templates(),
	}

	l, err := listConfig.Start()
	if err != nil {
		return err
	}

	var searchMode bool
	var searchInput []rune

	c.SetListener(func(line []rune, pos int, key rune) ([]rune, int, bool) {
		switch {
		case key == screen.KeyEnter:
			l.Open()
		case key == screen.KeyNext:
			l.Next()
		case key == screen.KeyPrev:
			l.Prev()
		case key == screen.KeySearch:
			searchMode = !searchMode
		case key == screen.KeyBackward:
			l.PageUp()
		case key == screen.KeyForward:
			l.PageDown()

		case key == screen.KeyBackspace:
			if searchMode {
				if len(searchInput) > 1 {
					searchInput = searchInput[:len(searchInput)-1]
				} else {
					l.Search(string(""))
					searchInput = nil
					searchMode = false
				}
			}
		default:
			if searchMode {
				searchInput = append(searchInput, line...)
				l.Search(string(searchInput))
			}
		}

		if searchMode {
			header := fmt.Sprintf(color.GreenString("Search: %s%s"), string(searchInput), screen.Cursor)
			_, err = t.WriteString(header)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_, err = t.WriteString(color.BlueString("Use the arrow keys to navigate: ↓ ↑ ← → / toggles search ↵ jumps to the build"))
			if err != nil {
				fmt.Println(err)
			}
		}

		items, active := l.Items()
		last := len(items) - 1

		for idx, item := range items {

			page := " "

			switch idx {
			case 0:
				if l.CanPageUp() {
					page = "↑"
				} else {
					page = " "
				}
			case last:
				if l.CanPageDown() {
					page = "↓"
				}
			}

			output := []byte(page + " ")
			if active == idx {
				output = append(output, l.Render(item, "Active")...)
			} else {
				output = append(output, l.Render(item, "Inactive")...)
			}

			_, err = t.Write(output)
			if err != nil {
				fmt.Println(err)
			}

		}

		if active == list.NotFound {
			_, err = t.Write([]byte("note found"))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			details := l.RenderDetails(items[active])
			for _, b := range details {
				_, err = t.Write(b)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		err = t.Flush()
		if err != nil {
			fmt.Println(err)
		}

		return nil, 0, true
	})

	for {
		_, err = rl.Readline()

		if err != nil {
			switch {
			case err == readline.ErrInterrupt, err.Error() == "Interrupt":
				fmt.Println(err)
			case err == io.EOF:
				fmt.Println(err)
			}
			break
		}
	}

	if err != nil {
		return err
	}

	return nil

}

func main() {
	flag.Parse()

	if helpBool {
		fmt.Printf("%v", desc)
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if versionBool {
		fmt.Printf("bkb v%v\n", version)
		os.Exit(0)
	}

	fmt.Printf("Error: %v\n", run())
}
