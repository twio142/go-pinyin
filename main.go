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
	fullWidthRegex = regexp.MustCompile(`[\p{Han}，。！？：；（）【】]+`)
	// Regex to match Pinyin characters (lowercase letters).
	pinyinRegex = regexp.MustCompile(`[a-z]+`)
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
	only     = flag.Bool("only", false, "Print only converted text, not original text")
	help     = flag.Bool("h", false, "Show usage message")
)

func init() {
	// Populate the full-width to half-width dictionary
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

	// Populate Xiaohe Shuangpin initial map
	xiaoheInitialMap = map[string]string{
		"b": "b", "p": "p", "m": "m", "f": "f",
		"d": "d", "t": "t", "n": "n", "l": "l",
		"g": "g", "k": "k", "h": "h",
		"j": "j", "q": "q", "x": "x",
		"z": "z", "c": "c", "s": "s",
		"zh": "v", "ch": "i", "sh": "u", "r": "r",
		"w": "w", "y": "y",
	}

	// Populate Xiaohe Shuangpin final map
	xiaoheFinalMap = map[string]string{
		"a": "a", "o": "o", "e": "e", "i": "i", "u": "u", "v": "u",
		"ai": "d", "ei": "w", "ao": "c", "ou": "z",
		"an": "j", "en": "f", "ang": "h", "eng": "g", "ong": "s", "iong": "s",
		"ia": "x", "ie": "p", "iao": "n", "iu": "q", "ian": "m", "in": "b", "iang": "l", "ing": "k",
		"ua": "x", "uo": "o", "uai": "k", "ui": "v", "uan": "r", "un": "y", "uang": "l",
		"ve": "t", "vn": "y",
	}

	// Populate Xiaohe Shuangpin zero-initial map (full syllable to two keys)
	xiaoheZeroInitialMap = map[string]string{
		"a": "aa", "e": "ee", "er": "er", "o": "oo",
	}
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -initials    Convert to Pinyin initials\n")
	fmt.Fprintf(os.Stderr, "  -xiaohe      Convert to Xiaohe Shuangpin\n")
	fmt.Fprintf(os.Stderr, "  -only        Print only converted text, not original text\n")
	fmt.Fprintf(os.Stderr, "  -h           Show this message\n")
	fmt.Fprintf(os.Stderr, "\nReads from standard input and converts Chinese characters to Pinyin.\n")
}

func main() {
	flag.Parse()

	// Check if help flag is set
	if *help {
		printUsage()
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		processLine(line, *initials, *xiaohe, *only)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func processLine(line string, isInitialMode bool, isXiaoheMode bool, isOnlyMode bool) {
	// In-place replace Chinese characters with Pinyin.
	result := fullWidthRegex.ReplaceAllStringFunc(line, func(s string) string {
		args := pinyin.NewArgs() // Args for normal Pinyin to get initials/finals
		args.Heteronym = true
		args.Fallback = func(r rune, a pinyin.Args) []string {
			if halfWidthChar, ok := fullWidthDict[r]; ok {
				return []string{string(halfWidthChar) + " "}
			}
			return []string{string(r)}
		}

		var stringBuilder strings.Builder

		if isXiaoheMode {
			for _, r := range s { // Iterate through each Hanzi character in the matched string 's'
				pinyinSyllables := pinyin.Pinyin(string(r), args) // Get Pinyin for single Hanzi
				if len(pinyinSyllables) > 0 {
					pinyinSyllables[0] = removeDuplicates(pinyinSyllables[0]) // Remove duplicates in the first slice
					for _, pinyinStr := range pinyinSyllables[0] {
						if !pinyinRegex.MatchString(pinyinStr) {
							stringBuilder.WriteString(pinyinStr)
						} else if mappedZeroInitial, ok := xiaoheZeroInitialMap[pinyinStr]; ok {
							if isInitialMode {
								stringBuilder.WriteString(mappedZeroInitial[:1])
							} else {
								stringBuilder.WriteString(mappedZeroInitial)
							}
						} else {
							initial := pinyinStr[0:2] // Get the first two characters as initial
							if mappedInitial, ok := xiaoheInitialMap[initial]; ok {
								stringBuilder.WriteString(mappedInitial)
								if !isInitialMode {
									final := pinyinStr[2:] // Get the rest as final
									if mappedFinal, ok := xiaoheFinalMap[final]; ok {
										stringBuilder.WriteString(mappedFinal)
									}
								}
							} else {
								initial := pinyinStr[0:1] // Fallback to first character if no mapping found
								if mappedInitial, ok := xiaoheInitialMap[initial]; ok {
									stringBuilder.WriteString(mappedInitial)
									if !isInitialMode {
										final := pinyinStr[1:] // Get the rest as final
										if mappedFinal, ok := xiaoheFinalMap[final]; ok {
											stringBuilder.WriteString(mappedFinal)
										}
									}
								} else {
									// Fallback to normal Pinyin if mapping not found (should not happen with complete maps)
									if isInitialMode {
										stringBuilder.WriteString(pinyinStr[:1])
									} else {
										stringBuilder.WriteString(pinyinStr)
									}
								}
							}
						}
						stringBuilder.WriteString(" ")
					}
				}
				// Add a space after each converted character's Xiaohe Pinyin
				stringBuilder.WriteString(" ")
			}
		} else {
			// Not Xiaohe mode
			if isInitialMode {
				args.Style = pinyin.FirstLetter
			} else {
				args.Style = pinyin.Normal // Default
			}

			for _, slice := range pinyin.Pinyin(s, args) {
				slice = removeDuplicates(slice)
				stringBuilder.WriteString(strings.Join(slice, " "))
				stringBuilder.WriteString(" ")
			}
		}
		return " " + strings.TrimSpace(stringBuilder.String()) + " "
	})

	// Clean up spacing.
	result = spaceCollapseRegex.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	// Handle output based on mode
	if isOnlyMode {
		fmt.Println(result)
	} else {
		// Default mode: print both lines if a conversion happened
		if result != line {
			fmt.Printf("%s\t%s\n", line, result)
		} else {
			fmt.Println(line)
		}
	}
}
