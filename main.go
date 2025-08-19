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

	// Xiaohe Shuangpin mapping tables
	xiaoheInitialMap     = map[string]string{}
	xiaoheFinalMap       = map[string]string{}
	xiaoheZeroInitialMap = map[string]string{}

	// Command-line flags
	initials = flag.Bool("initials", false, "Convert to Pinyin initials")
	xiaohe   = flag.Bool("xiaohe", false, "Convert to Xiaohe Shuangpin")
)

func init() {
	// Populate the full-width to half-width dictionary
	for i := range 26 {
		fullWidthDict[rune(0xFF21+i)] = rune('A' + i) // A-Z
		fullWidthDict[rune(0xFF41+i)] = rune('a' + i) // a-z
	}
	for i := range 10 {
		fullWidthDict[rune(0xFF10+i)] = rune('0' + i) // 0-9
	}
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

	// Populate Xiaohe Shuangpin initial map
	xiaoheInitialMap = map[string]string{
		"b": "b", "p": "p", "m": "m", "f": "f",
		"d": "d", "t": "t", "n": "n", "l": "l",
		"g": "g", "k": "k", "h": "h",
		"j": "j", "q": "q", "x": "x",
		"z": "z", "c": "c", "s": "s",
		"zh": "v", "ch": "i", "sh": "u", "r": "r",
	}

	// Populate Xiaohe Shuangpin final map
	xiaoheFinalMap = map[string]string{
		"a": "a", "o": "o", "e": "e", "i": "i", "u": "u", "v": "u",
		"ai": "d", "ei": "w", "ao": "k", "ou": "z",
		"an": "j", "en": "f", "ang": "h", "eng": "g", "ong": "s", "iong": "s",
		"ia": "x", "ie": "p", "iao": "n", "iu": "q", "ian": "m", "in": "b", "iang": "l", "ing": "y",
		"ua": "x", "uo": "o", "uai": "k", "ui": "v", "uan": "r", "un": "y", "uang": "l",
		"ve": "t", "vn": "y",
	}

	// Populate Xiaohe Shuangpin zero-initial map (full syllable to two keys)
	xiaoheZeroInitialMap = map[string]string{
		"a": "aa", "e": "ee", "er": "er", "o": "oo",
	}
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		processLine(line, a, *initials, *xiaohe)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func processLine(line string, args pinyin.Args, isInitialMode bool, isXiaoheMode bool) {
	// Convert full-width characters to half-width using the dictionary.
	halfWidthLine := convertFullWidth(line)

	// In-place replace Chinese characters with Pinyin.
	result := hanRegex.ReplaceAllStringFunc(halfWidthLine, func(s string) string {
		var convertedPinyin string

		if isXiaoheMode {
			var xiaoheBuilder strings.Builder
			normalArgs := pinyin.NewArgs() // Args for normal Pinyin to get initials/finals

			for _, r := range s { // Iterate through each Hanzi character in the matched string 's'
				pinyinSyllables := pinyin.Pinyin(string(r), normalArgs) // Get Pinyin for single Hanzi
				if len(pinyinSyllables) > 0 && len(pinyinSyllables[0]) > 0 {
					pinyinStr := pinyinSyllables[0][0] // Get the first (normal) Pinyin pronunciation

					// Check for zero-initials first
					if mappedZeroInitial, ok := xiaoheZeroInitialMap[pinyinStr]; ok {
						xiaoheBuilder.WriteString(mappedZeroInitial)
					} else {
						// Not a zero-initial, get initial and final
						initialsForRune := pinyin.Pinyin(string(r), pinyin.Args{Style: pinyin.Initials})
						finalsForRune := pinyin.Pinyin(string(r), pinyin.Args{Style: pinyin.Finals})

						initial := ""
						final := ""

						if len(initialsForRune) > 0 && len(initialsForRune[0]) > 0 {
							initial = initialsForRune[0][0]
						}
						if len(finalsForRune) > 0 && len(finalsForRune[0]) > 0 {
							final = finalsForRune[0][0]
						}

						mappedInitial, initialExists := xiaoheInitialMap[initial]
						mappedFinal, finalExists := xiaoheFinalMap[final]

						if initialExists && finalExists {
							xiaoheBuilder.WriteString(mappedInitial)
							xiaoheBuilder.WriteString(mappedFinal)
						} else {
							// Fallback to normal Pinyin if mapping not found (should not happen with complete maps)
							xiaoheBuilder.WriteString(pinyinStr)
						}
					}
				}
				// Add a space after each converted character's Xiaohe Pinyin
				xiaoheBuilder.WriteString(" ")
			}
			convertedPinyin = strings.TrimSpace(xiaoheBuilder.String())

			if isInitialMode {
				parts := strings.Split(convertedPinyin, " ")
				initialsSlice := make([]string, 0, len(parts))
				for _, part := range parts {
					if len(part) > 0 {
						initialsSlice = append(initialsSlice, string(part[0]))
					}
				}
				convertedPinyin = strings.Join(initialsSlice, " ")
			}
		} else {
			// Not Xiaohe mode, check for initials or normal
			if isInitialMode {
				args.Style = pinyin.FirstLetter
			} else {
				args.Style = pinyin.Normal // Default
			}
			pinyinSlices := pinyin.LazyConvert(s, &args)
			convertedPinyin = strings.Join(pinyinSlices, " ")
		}

		// Pad with spaces to ensure separation from non-Chinese parts.
		return " " + convertedPinyin + " "
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
