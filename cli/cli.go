package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"hash/crc64"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/signal"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/uc/cli/support"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application basic info
const (
	APP  = "uc"
	VER  = "2.0.1"
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

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MAX_SAMPLE_SIZE is maximum sample size
const MAX_SAMPLE_SIZE = 512

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
	OPT_HELP:         {Type: options.BOOL},
	OPT_VER:          {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// stats contains info about data
var stats *Stats

// rawMode is raw mode flag
var rawMode bool

// colorTagApp contains color tag for app name
var colorTagApp string

// colorTagVer contains color tag for app version
var colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main application function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(1)

	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	signal.Handlers{
		signal.INT:  signalHandler,
		signal.TERM: signalHandler,
		signal.QUIT: signalHandler,
	}.TrackAsync()

	processData(args)
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
		rawMode = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_NO_PROGRESS) {
		rawMode = true
	}

	switch {
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#168}", "{#168}"
	default:
		colorTagApp, colorTagVer = "{*}{m}", "{m}"
	}
}

// processData starts data processing
func processData(args options.Arguments) {
	var r *bufio.Reader

	stats = &Stats{
		Counters: make(map[uint64]uint32),
		mx:       &sync.Mutex{},
	}

	input := getInput(args)

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

// getInput returns input for reading data
func getInput(args options.Arguments) string {
	if args.Get(0).String() == "-" || !fsutil.IsCharacterDevice("/dev/stdin") {
		return "-"
	}

	input := args.Get(0).Clean().String()
	err := fsutil.ValidatePerms("FRS", input)

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	return input
}

// readData reads data
func readData(s *bufio.Scanner) {
	ct := crc64.MakeTable(crc64.ECMA)
	dist := options.GetB(OPT_DISTRIBUTION)
	maxLines, err := parseMaxLines(options.GetS(OPT_MAX_LINES))

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

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
				stats.Samples[dataCrc] = strutil.Substr(string(data), 0, MAX_SAMPLE_SIZE)
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
	for range time.NewTicker(time.Second / 4).C {
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

// printDistribution prints distribution info
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

// parseMaxLines parses max line option
func parseMaxLines(maxLines string) (int, error) {
	if maxLines == "" {
		return 0, nil
	}

	maxLines = strings.ToUpper(maxLines)

	mp := 1

	switch {
	case strings.HasSuffix(maxLines, "K"):
		maxLines = strutil.Exclude(maxLines, "K")
		mp = 1000
	case strings.HasSuffix(maxLines, "M"):
		mp = 1000 * 1000
		maxLines = strutil.Exclude(maxLines, "M")
	}

	num, err := strconv.Atoi(maxLines)

	if err != nil {
		return 0, err
	}

	return num * mp, nil
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

// printCompletion prints completion for given shell
func printCompletion() int {
	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(genUsage(), APP))
	case "fish":
		fmt.Print(fish.Generate(genUsage(), APP))
	case "zsh":
		fmt.Print(zsh.Generate(genUsage(), optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo(APP, "file")

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_DISTRIBUTION, "Show number of occurrences for every line")
	info.AddOption(OPT_MAX_LINES, "Max number of unique lines", "num")
	info.AddOption(OPT_NO_PROGRESS, "Disable progress output")
	info.AddOption(OPT_NO_PROGRESS, "Disable progress output")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("file.txt", "Count unique lines in file.txt")
	info.AddExample("-d file.txt", "Show distribution for file.txt")
	info.AddExample("-d -m 5k file.txt", "Show distribution for file.txt with 5,000 uniq lines max")
	info.AddRawExample("cat file.txt | "+APP, "Count unique lines in stdin data")

	info.AddRawExample(
		APP+" -m 100 < file.txt",
		"Count unique lines in stdin data with 100 uniq lines max",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		BugTracker:    "https://github.com/essentialkaos/uc",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/uc", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
