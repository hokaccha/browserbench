package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/sourcegraph/go-selenium"
)

type BenchmarkOpts struct {
	Url     string
	Start   string
	End     string
	Number  int
	Browser string
	Remote  string
}

type PerformanceEntry struct {
	Duration  float64 `json:duration`
	EntryType string  `json:entryType`
	Name      string  `json:name`
	StartTime float64 `json:startTime`
}

var Result = []PerformanceEntry{}

const jsTmpl = `
	console.log('foo');
	var done = arguments[arguments.length - 1];
	(function loop() {
		if (window.performance.timing['{{.End}}'] === 0 && window.performance.getEntriesByName('{{.End}}').length === 0) {
			return setTimeout(loop, 200);
		}
		window.performance.measure('{{.Name}}', '{{.Start}}', '{{.End}}');
		done(window.performance.getEntriesByName('{{.Name}}')[0]);
	})();
`

func Benchmark(opts *BenchmarkOpts) {
	for i := 0; i < opts.Number; i++ {
		bench(opts)
	}

	fmt.Printf("\nAverage: %dms\n", int(calcAverage()))
}

func bench(opts *BenchmarkOpts) {
	caps := selenium.Capabilities(map[string]interface{}{"browserName": opts.Browser})
	wd, err := selenium.NewRemote(caps, opts.Remote)

	if err != nil {
		log.Fatal(err)
	}

	defer wd.Quit()

	err = wd.Get(opts.Url)

	if err != nil {
		log.Fatal(err)
	}

	entry := getEntry(wd, opts)

	fmt.Printf("%s: %dms\n", entry.Name, int(entry.Duration))
	Result = append(Result, *entry)
}

func getEntry(wd selenium.WebDriver, opts *BenchmarkOpts) *PerformanceEntry {
	err := wd.SetAsyncScriptTimeout(20000)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("js").Parse(jsTmpl)

	if err != nil {
		log.Fatal(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, map[string]string{
		"Name":  opts.Start + "..." + opts.End,
		"Start": opts.Start,
		"End":   opts.End,
	})

	if err != nil {
		log.Fatal(err)
	}

	js := doc.String()
	result, err := wd.ExecuteScriptAsync(js, nil)

	if err != nil {
		log.Fatal(err)
	}

	performanceEntry := &PerformanceEntry{}
	byt, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byt, performanceEntry)

	return performanceEntry
}

func calcAverage() float64 {
	var sum float64
	for _, entry := range Result {
		sum += entry.Duration
	}
	return sum / float64(len(Result))
}
