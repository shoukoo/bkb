package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/shoukoo/bkb/list"
	"github.com/shoukoo/bkb/screen"
)

type tpe struct {
	Name string
	Type int
}

func run() error {
	// Readline reads input from user
	// TODO not fully understand how this pkg works
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

	rl.Write([]byte(screen.HideCursor))
	t := screen.New(rl)

	client, err := list.BuildkiteClient()
	if err != nil {
		return err
	}
	client.GetRecentBuilds("lexer")

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
			header := fmt.Sprintf("Search: %s%s", string(searchInput), screen.Cursor)
			t.WriteString(header)
		} else {
			t.WriteString("Use the arrow keys to navigate: ↓ ↑ ← → / toggles search ↵ jump to the build")
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

			t.Write(output)

		}

		if active == list.NotFound {
			t.Write([]byte("note found"))
		} else {
			details := l.RenderDetails(items[active])
			for _, b := range details {
				t.Write(b)
			}
		}

		t.Flush()

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
	fmt.Printf("Error: %v", run())
}
