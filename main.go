package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

var (
	// Regex to detect consecutive Chinese characters.
	hanRegex = regexp.MustCompile(`[\p{Han}]+`)
	// Regex to collapse multiple spaces.
	spaceCollapseRegex = regexp.MustCompile(`\s{2,}`)

	// Dictionary for full-width to half-width conversion.
	fullWidthDict = map[rune]rune{}

	// Command-line flags
	initials = flag.Bool("initials", false, "Convert to Pinyin initials")
)

func init() {
	// Populate the dictionary programmatically.
	// Full-width latin letters (A-Z, a-z)
	for i := 0; i < 26; i++ {
		fullWidthDict[rune(0xFF21+i)] = rune('A'+i)
		fullWidthDict[rune(0xFF41+i)] = rune('a'+i)
	}
	// Full-width numbers (0-9)
	for i := 0; i < 10; i++ {
		fullWidthDict[rune(0xFF10+i)] = rune('0'+i)
	}
	// Full-width punctuation
	fullWidthDict['，'] = ','
	fullWidthDict['。'] = '.'
	fullWidthDict['？'] = '?'
	fullWidthDict['！'] = '!'
	fullWidthDict['；'] = ';'
	fullWidthDict['：'] = ':'
	fullWidthDict['（'] = '('
	fullWidthDict['）'] = ')'
	fullWidthDict['【'] = '['
	fullWidthDict['】'] = ']'
	fullWidthDict['“'] = '"'
	fullWidthDict['”'] = '"'
	fullWidthDict['‘'] = '\''
	fullWidthDict['’'] = '\''
	fullWidthDict['　'] = ' '
}

func convertFullWidth(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if replacement, ok := fullWidthDict[r]; ok {
			b.WriteRune(replacement)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func main() {
	flag.Parse()

	// Configure pinyin conversion arguments based on flags.
	a := pinyin.NewArgs()
	if *initials {
		a.Style = pinyin.Initials
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		processLine(line, a, *initials)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func processLine(line string, args pinyin.Args, isInitialMode bool) {
	// Convert full-width characters to half-width using the dictionary.
	halfWidthLine := convertFullWidth(line)

	// In-place replace Chinese characters with Pinyin.
	result := hanRegex.ReplaceAllStringFunc(halfWidthLine, func(s string) string {
		pinyinSlices := pinyin.LazyConvert(s, &args)
		joiner := " "
		// Initials are joined without spaces.
		if isInitialMode {
			joiner = ""
		}
		// Pad with spaces to ensure separation from non-Chinese parts.
		return " " + strings.Join(pinyinSlices, joiner) + " "
	})

	// Clean up spacing.
	result = spaceCollapseRegex.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	// Only print both lines if a conversion happened.
	if result != line {
		fmt.Printf("%s\t%s\n", line, result)
	} else {
		fmt.Println(line)
	}
}
