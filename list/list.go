package list

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

const NotFound = -1

type Config struct {
	Items     interface{}
	Templates *SelectTemplates
}

func (c *Config) Start() (*List, error) {
	l, err := New(c.Items, 5)
	if err != nil {
		return nil, err
	}

	l.Templates = c.Templates
	err = l.prepareTemplates()
	if err != nil {
		return nil, err
	}

	return l, nil
}

type List struct {
	items     []interface{}
	scope     []interface{} // list to display on the screen
	cursor    int
	size      int
	start     int
	Templates *SelectTemplates
}

func (s *List) prepareTemplates() error {
	tpls := s.Templates
	if tpls == nil {
		tpls = &SelectTemplates{}
	}

	tpl, err := template.New("").Funcs(FuncMap).Parse(tpls.Active)
	if err != nil {
		return err
	}
	tpls.active = tpl

	tpl, err = template.New("").Funcs(FuncMap).Parse(tpls.Inactive)
	if err != nil {
		return err
	}
	tpls.inactive = tpl

	tpl, err = template.New("").Funcs(FuncMap).Parse(tpls.Details)
	if err != nil {
		return err
	}

	tpls.details = tpl

	s.Templates = tpls

	return nil
}

func New(item interface{}, size int) (*List, error) {

	if size < 1 {
		return nil, fmt.Errorf("list size %d must be greater than 0", size)
	}

	if item == nil {
		return nil, fmt.Errorf("selection is empty")
	}

	if reflect.TypeOf(item).Kind() != reflect.Slice {
		return nil, fmt.Errorf("item %v is not a slice", item)
	}

	slice := reflect.ValueOf(item)
	values := make([]interface{}, slice.Len())

	for i := range values {
		v := slice.Index(i).Interface()
		values[i] = v
	}

	return &List{size: size, items: values, scope: values}, nil

}

func (l *List) Items() ([]interface{}, int) {
	var result []interface{}
	max := len(l.scope)     //total no. of list
	end := l.start + l.size //max size can display on the screen

	if end > max {
		end = max
	}

	active := NotFound

	// i = l.start and j = 0. j is the index of the loop.
	// l.start < end the continue
	// in the end i+1 and j+1
	for i, j := l.start, 0; i < end; i, j = i+1, j+1 {

		// if cursore == i (l.start)
		if l.cursor == i {
			active = j
		}

		result = append(result, l.scope[i])
	}

	return result, active

}

func (l *List) Open() {
	v := reflect.ValueOf(l.scope[l.cursor])

	// if its a pointer, resolve its value
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	if v.Kind() != reflect.Struct {
		log.Fatalf("Not an interface %v", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {

		if v.Type().Field(i).Name == "WebURL" {

			err := open(v.Field(i).String())
			if err != nil {
				log.Fatalf("%v", err)
			}
			return

		}

	}

	log.Fatalf("Can't find URL")

}

func (l *List) Next() {

	max := len(l.scope) - 1

	if l.cursor < max {
		l.cursor++
	}

	// if current display list is smaller than cursor
	if l.start+l.size <= l.cursor {
		l.start = l.cursor - l.size + 1
	}
}

func (l *List) Prev() {

	// cursor needs to go to the top
	// then change start accordingly
	if l.cursor > 0 {
		l.cursor--
	}

	if l.start > l.cursor {
		l.start = l.cursor
	}
}

// PageUp moves the visible list backward by x items. Where x is the size of the
// visible items on the list. The selected item becomes the first visible item.
// If the list is already at the bottom, the selected item becomes the last
// visible item.
func (l *List) PageUp() {
	start := l.start - l.size
	if start < 0 {
		l.start = 0
	} else {
		l.start = start
	}

	cursor := l.start

	if cursor < l.cursor {
		l.cursor = cursor
	}
}

// PageDown moves the visible list forward by x items. Where x is the size of
// the visible items on the list. The selected item becomes the first visible
// item.
func (l *List) PageDown() {
	start := l.start + l.size
	max := len(l.scope) - l.size

	switch {
	case len(l.scope) < l.size:
		l.start = 0
	case start > max:
		l.start = max
	default:
		l.start = start
	}

	cursor := l.start

	if cursor == l.cursor {
		l.cursor = len(l.scope) - 1
	} else if cursor > l.cursor {
		l.cursor = cursor
	}
}

// CanPageDown returns whether a list can still PageDown().
func (l *List) CanPageDown() bool {
	max := len(l.scope)
	return l.start+l.size < max
}

// CanPageUp returns whether a list can still PageUp().
func (l *List) CanPageUp() bool {
	return l.start > 0
}

// Search
func (l *List) Search(key string) {

	key = strings.Trim(key, " ")
	l.cursor = 0
	l.start = 0
	l.search(key)
}

func (l *List) search(key string) {
	if len(key) == 0 {
		l.scope = l.items
		return
	}
	var scope []interface{}
	for _, v := range l.items {
		if ok := strings.Contains(fmt.Sprint(v), key); ok {
			scope = append(scope, v)
		}
	}

	l.scope = scope
}

func (l *List) Render(item interface{}, tpe string) []byte {
	var buf bytes.Buffer
	switch tpe {
	case "Active":
		err := l.Templates.active.Execute(&buf, item)
		if err != nil {
			buf.WriteString(err.Error())
		}
	case "Inactive":
		err := l.Templates.inactive.Execute(&buf, item)
		if err != nil {
			buf.WriteString(err.Error())
		}
	default:
		buf.WriteString(fmt.Sprintf("No template found for %s", tpe))
	}

	return buf.Bytes()
}

func (l *List) RenderDetails(item interface{}) [][]byte {
	var buf bytes.Buffer
	err := l.Templates.details.Execute(&buf, item)
	if err != nil {
		buf.WriteString(err.Error())
	}

	output := buf.Bytes()
	return bytes.Split(output, []byte("\n"))
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
