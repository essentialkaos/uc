package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"hash/crc64"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"pkg.re/essentialkaos/ek.v11/fmtc"
	"pkg.re/essentialkaos/ek.v11/fmtutil"
	"pkg.re/essentialkaos/ek.v11/fsutil"
	"pkg.re/essentialkaos/ek.v11/options"
	"pkg.re/essentialkaos/ek.v11/signal"
	"pkg.re/essentialkaos/ek.v11/strutil"
	"pkg.re/essentialkaos/ek.v11/usage"
	"pkg.re/essentialkaos/ek.v11/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v11/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v11/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v11/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application basic info
const (
	APP  = "uc"
	VER  = "0.0.2"
	DESC = "Tool for counting unique lines"
)

// Constants with options names
const (
	OPT_MAX_LINES    = "m:max"
	OPT_DISTRIBUTION = "d:dist"
	OPT_NO_PROGRESS  = "np:no-progress"
	OPT_NO_COLOR     = "nc:no-color"
	OPT_HELP         = "h:help"
	OPT_VER          = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Stats contains data info
type Stats struct {
	Counters       map[uint64]uint32 // crc64 → num
	Samples        map[uint64]string // crc64 → sample (512 symbols)
	LastReadLines  uint64
	LastReadBytes  float64
	TotalReadLines uint64
	TotalReadBytes float64
	LastReadDate   time.Time
	Finished       bool

	mx *sync.Mutex
}

// ////////////////////////////////////////////////////////////////////////////////// //

// LineInfo is struct with line info
type LineInfo struct {
	CRC uint64
	Num uint32
}

type linesSlice []LineInfo

func (s linesSlice) Len() int      { return len(s) }
func (s linesSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s linesSlice) Less(i, j int) bool {
	return s[i].Num < s[j].Num
}

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap is map with options
var optMap = options.Map{
	OPT_MAX_LINES:    {Type: options.INT},
	OPT_DISTRIBUTION: {Type: options.BOOL},
	OPT_NO_PROGRESS:  {Type: options.BOOL},
	OPT_NO_COLOR:     {Type: options.BOOL},
	OPT_HELP:         {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:          {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// stats contains info about data
var stats *Stats

// rawMode is raw mode flag
var rawMode bool

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main func
func main() {
	runtime.GOMAXPROCS(1)

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError("Options parsing errors:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		os.Exit(0)
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		os.Exit(0)
	}

	signal.Handlers{
		signal.INT:  signalHandler,
		signal.TERM: signalHandler,
		signal.QUIT: signalHandler,
	}.TrackAsync()

	processData(args[0])
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if !fsutil.IsCharacterDevice("/dev/stdout") && os.Getenv("FAKETTY") == "" {
		rawMode = true
	}

	if options.GetB(OPT_NO_PROGRESS) {
		rawMode = true
	}
}

// processData starts data processing
func processData(input string) {
	var r *bufio.Reader

	stats = &Stats{
		Counters: make(map[uint64]uint32),
		mx:       &sync.Mutex{},
	}

	if input == "-" {
		r = bufio.NewReader(os.Stdin)
	} else {
		fd, err := os.OpenFile(input, os.O_RDONLY, 0)

		if err != nil {
			printError(err.Error())
			os.Exit(1)
		}

		r = bufio.NewReader(fd)
	}

	readData(bufio.NewScanner(r))
}

// readData reads data
func readData(s *bufio.Scanner) {
	ct := crc64.MakeTable(crc64.ECMA)
	dist := options.GetB(OPT_DISTRIBUTION)
	maxLines := options.GetI(OPT_MAX_LINES)

	if dist {
		stats.Samples = make(map[uint64]string)
	}

	stats.LastReadDate = time.Now()

	if !rawMode {
		go printProgress()
	}

	for s.Scan() {
		data := s.Bytes()
		dataLen := float64(len(data))
		dataCrc := crc64.Checksum(data, ct)

		stats.mx.Lock()

		stats.Counters[dataCrc]++
		stats.LastReadBytes += dataLen
		stats.LastReadLines++

		stats.TotalReadLines++
		stats.TotalReadBytes += dataLen

		if dist {
			_, exist := stats.Samples[dataCrc]

			if !exist {
				stats.Samples[dataCrc] = strutil.Substr(string(data), 0, 512)
			}
		}

		if maxLines > 0 && len(stats.Counters) == maxLines {
			stats.mx.Unlock()
			break
		}

		stats.mx.Unlock()
	}

	printResults()
}

// printProgress shows data processing progress
func printProgress() {
	for range time.NewTicker(time.Second / 5).C {
		stats.mx.Lock()

		if stats.Finished {
			break
		}

		now := time.Now()
		dur := now.Sub(stats.LastReadDate)
		readSpeed := stats.LastReadBytes / dur.Seconds()

		fmtc.TPrintf(
			"{s}%12s/s {s-}|{s} %-12s {s-}|{s} %12s/s {s-}|{s} %-12s{!}",
			fmtutil.PrettyNum(stats.LastReadLines),
			fmtutil.PrettyNum(stats.TotalReadLines),
			fmtutil.PrettySize(readSpeed),
			fmtutil.PrettySize(stats.TotalReadBytes),
		)

		stats.LastReadLines = 0
		stats.LastReadBytes = 0
		stats.LastReadDate = now

		stats.mx.Unlock()
	}
}

// printResults shows results
func printResults() {
	stats.mx.Lock()

	stats.Finished = true

	if options.GetB(OPT_DISTRIBUTION) {
		printDistribution()
	} else {
		fmtc.TPrintln(len(stats.Counters))
	}

	stats.mx.Unlock()
}

// printDistribution prints distrubution info
func printDistribution() {
	var distData linesSlice

	for crc, num := range stats.Counters {
		distData = append(distData, LineInfo{crc, num})
	}

	sort.Sort(sort.Reverse(distData))

	for _, info := range distData {
		fmtc.TPrintf(" %7d %s\n", info.Num, stats.Samples[info.CRC])
	}
}

// signalHandler is signal handler
func signalHandler() {
	printResults()
	os.Exit(0)
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage print usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo(APP, "file")

	info.AddOption(OPT_DISTRIBUTION, "Show number of occurrences for every line")
	info.AddOption(OPT_MAX_LINES, "Max number of unique lines {s-}(default: 5000){!}", "num")
	info.AddOption(OPT_NO_PROGRESS, "Disable progress output")
	info.AddOption(OPT_NO_PROGRESS, "Disable progress output")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("file.txt", "Count unique lines in file.txt")
	info.AddExample("-d file.txt", "Show distribution for file.txt")
	info.AddRawExample(
		"cat file.txt | "+APP+" -",
		"Count unique lines in stdin data",
	)

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(genUsage(), APP))
	case "fish":
		fmt.Printf(fish.Generate(genUsage(), APP))
	case "zsh":
		fmt.Printf(zsh.Generate(genUsage(), optMap, APP))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout print info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "Essential Kaos",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/uc", update.GitHubChecker},
	}

	about.Render()
}
