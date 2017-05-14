package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// create crossword datastructure
//  (borrowed from tic-tac-toe in golang tour
const cw_dim = 15       // crosswords is 15 across
const len_alphabet = 26 // length of the alphabet
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var cw [cw_dim][cw_dim]string         // crossword 2D array, everything is a string
var ref_box = make(map[string]string) // reference box, no numbers all strings
var wordlist [cw_dim + 1][]string     // wordlist  mapped on word length (array starts at 0)

func main() {
	ReadCrossword("crossword-7") // read this crossword and populate data structure
	InitialiseRefbox()
	// look for all words
	words := GetWordsCW()

	// read the wordlist
	ReadWordlist("english-words/words.txt")

	//fmt.Println(wordlist)
	//fmt.Println(wordlist[4])

	PrintCw(words)
	for i := 1; i <= 15; i++ {
		//		// print the current CrossWord and the Refbox
		fmt.Printf("============\nIteration: %d\n=============\n", i)
		PrintRefbox()
		change := TryWordlist(words)

		if !change {
			break
		}
		//
		//		// make a choice
		//		MakeChoiceRefbox()
	}
	//
	PrintRefbox()
	PrintCw(words)
}

func TryWordlist(words [][]string) bool {
	// try matchiong wordlist entries to the crossword
	change := false                        // is there a change in this iteration
	re_char := regexp.MustCompile(`[A-Z]`) // upper case Char
	max_hits := 5                          // the max number of hits we will present to use

	_, incomplete := AnalyseCw(words)
	unfound := UnfoundChars()
	r_unfound := "[" // char class to be used in regexp for all non-found
	for _, u := range unfound {
		r_unfound += u
	}
	r_unfound += "]"
	//fmt.Println("unfound: ", unfound, r_unfound)

	for _, word := range incomplete {
		var hits []string
		var regex string
		var p_word string // string representation of word to print (without regex)

		for _, c := range word {
			// checken tov ref_box[c] !!!
			if re_char.MatchString(ref_box[c]) {
				regex += ref_box[c]
				p_word += ref_box[c]
			} else {
				regex += r_unfound
				p_word += "."
			}
		}
		fmt.Println("tlw word - regex: ", word, p_word, regex)
		if p_word == regex { // this word became complete by another word in this iteration
			continue
		}
		re_word := regexp.MustCompile(regex)
		num_hits := 0
		for _, wl := range wordlist[len(word)] {

			if re_word.MatchString(wl) {
				if num_hits == max_hits { // after max_hits we stop
					num_hits++ // to indicate more then max_hits
					break
				}
				hits = append(hits, wl)
				num_hits++
			}
		}
		if num_hits >= max_hits { // too many hits, move on
			continue
		}
		fmt.Println("   match with ", hits)

		if num_hits == 1 { // only a single match, update ref_box
			change = true
			letters := strings.Split(hits[0], "")
			for i := 0; i < len(word); i++ { // set everything in ref_box
				ref_box[word[i]] = letters[i]
			}
			continue
		}

		// we have 2 - max_hits hits, find the common chars in all of them
		for i := 0; i < len(word); i++ {
			h_char := make(map[string]int) // hash to store the different chars on this position
			var last_char string           // track the last char we saw
			for _, hit := range hits {
				letters := strings.Split(hit, "")
				h_char[letters[i]]++
				last_char = letters[i]
			}
			if h_char[last_char] == num_hits { // on this position we have the same char in all hits
				change = true
				fmt.Println("  i", i, " t", last_char, " ", h_char[last_char], word[i], ref_box[word[i]])
				ref_box[word[i]] = last_char
			}

		}
	}
	return change
}

func UnfoundChars() []string { // return a hash with the unfound chars
	letters := strings.Split(alphabet, "")
	h := make(map[string]bool) // hash to record refbox entries die er al zijn
	var unfound []string

	for i := 1; i <= len_alphabet; i++ { //  and the entries of ref_box
		if strconv.Itoa(i) != ref_box[strconv.Itoa(i)] { // we have a letter
			h[ref_box[strconv.Itoa(i)]] = true
		}
	}

	for _, l := range letters {
		if _, ok := h[l]; !ok {
			unfound = append(unfound, l)
		}
	}
	return unfound
}

func PrintCw(words [][]string) {
	complete, incomplete := AnalyseCw(words)

	for _, word := range complete {
		fmt.Printf("\t")
		for _, c := range word {
			fmt.Printf("%s", ref_box[c])
		}
		fmt.Println()
	}
	for _, word := range incomplete {
		fmt.Print(" ")
		for _, c := range word {
			fmt.Printf("%2s ", ref_box[c])
		}
		fmt.Println()
	}
}

func AnalyseCw(words [][]string) ([][]string, [][]string) {
	var complete, incomplete [][]string
	for _, word := range words {
		incomp := true // by default incomplete
		for _, c := range word {
			//fmt.Println("c", c, "ref_b", ref_box[c])
			if c == ref_box[c] {
				incomp = true // still a number in the word
				break
			} else {
				incomp = false //   a char
			}
		}
		if incomp == true {
			incomplete = append(incomplete, word)
		} else { // word already completed
			complete = append(complete, word)
		}
	}
	return complete, incomplete
}

func PrintRefbox() {
	letters := strings.Split(alphabet, "")
	h := make(map[string]bool) // hash to record refbox entries die er al zijn

	fmt.Printf("--------------------\n")
	for i := 1; i <= len_alphabet; i++ { // print 1 - 26
		fmt.Printf("%2d ", i)
	}
	fmt.Println()
	for i := 1; i <= len_alphabet; i++ { //  and the entries of ref_box
		if strconv.Itoa(i) != ref_box[strconv.Itoa(i)] { // we have a letter
			h[ref_box[strconv.Itoa(i)]] = true
			fmt.Printf("%2s ", ref_box[strconv.Itoa(i)])
		} else {
			fmt.Print("   ")
		}
	}
	fmt.Printf("\n--------------------\n")

	for _, l := range letters {
		if _, ok := h[l]; ok {
			fmt.Print("  ")
		} else {
			fmt.Printf("%2s", l)
		}
	}
	fmt.Println()
}

func MakeChoiceRefbox() {
	// make a choice in the refbox
	reader := bufio.NewReader(os.Stdin)
	re_mul_spaces := regexp.MustCompile(`\s+`) // replace multiple spaces by a single

	fmt.Println("----------------------------------")
	fmt.Print("Make choice for Refbox (num char): ")
	choice, _ := reader.ReadString('\n')
	choice = re_mul_spaces.ReplaceAllString(choice, " ")
	choice = strings.TrimSpace(choice)
	c := strings.Split(choice, " ")
	for k, v := range c {
		//fmt.Println("k", k, "v", v)
		ref_box[v] = c[k+1] // num, char
		break
	}

}

func GetWordsCW() [][]string {
	// look for all words (2 or more consecutive chars) in cw
	var words [][]string
	// first look in rows (easy)
	for i := 0; i < cw_dim; i++ {
		var word []string
		in_word := false
		for j := 0; j < cw_dim; j++ {
			if cw[i][j] == "*" {
				if in_word == true { // a * and in a word, append word[] to words[]
					words = append(words, word)
					word = nil // re-initialise word[]
					//word = word[:0] // re-initialise word[]
					in_word = false
				} // if not in_word, nothing to be done
			} else { // we have a number
				if in_word == false { // first nubmer
					if j+1 < cw_dim && cw[i][j+1] != "*" { // not the last cell and next cell is number
						in_word = true
						word = append(word, cw[i][j])
					}
				} else { // already in_word , add content to word[]
					word = append(word, cw[i][j])
				}
			}
		}
		// if word exists, add to words !! XXX
		if in_word == true {
			words = append(words, word)
		}
	}
	// then look at columns -- horrible duplication of code
	//  swap i and J
	for j := 0; j < cw_dim; j++ {
		var word []string
		in_word := false
		for i := 0; i < cw_dim; i++ {
			if cw[i][j] == "*" {
				if in_word == true { // a * and in a word, append word[] to words[]
					words = append(words, word)
					word = nil // re-initialise word[]
					in_word = false
				} // if not in_word, nothing to be done
			} else { // we have a number
				if in_word == false { // first nubmer
					if i+1 < cw_dim && cw[i+1][j] != "*" { // not the last cell and next cell is number
						in_word = true
						word = append(word, cw[i][j])
					}
				} else { // already in_word , add content to word[]
					word = append(word, cw[i][j])
				}
			}
		}
		// (end of column) if word exists, add to words !! XXX
		if in_word == true {
			words = append(words, word)
		}
	}
	return words
}

func InitialiseRefbox() {
	for i := 1; i <= len_alphabet; i++ {
		ref_box[strconv.Itoa(i)] = strconv.Itoa(i)
	}
	// We get three letters
	ref_box["2"] = "T"
	ref_box["10"] = "E"
	ref_box["12"] = "N"
}

func ReadCrossword(filepath string) {
	lines := File2Lines(filepath)

	re_lt_spaces := regexp.MustCompile(`^\s*(\S.*\S)\s*$`) // trailing and leading spaces
	re_mul_spaces := regexp.MustCompile(`\s+`)             // replace multiple spaces by a single
	re_mul_stars := regexp.MustCompile(`\*\*`)             // replace multiple stars by a single

	var j = 0
	for _, line := range lines {
		line = re_lt_spaces.ReplaceAllString(line, "$1")
		line = re_mul_spaces.ReplaceAllString(line, " ")
		line = re_mul_stars.ReplaceAllString(line, "*")
		cell := strings.Split(line, " ")
		var k = 0
		for _, c := range cell { // populate cw 2D array
			cw[j][k] = c
			k++
		}
		j++
	}
}

func ReadWordlist(filepath string) {
	fmt.Print("Reading wordlist ...")
	lines := File2Lines(filepath)
	fmt.Print(" read file ... ")

	//re_lt_spaces := regexp.MustCompile(`^\s*(\S.*\S)\s*$`) // trailing and leading spaces
	//re_mul_spaces := regexp.MustCompile(`\s+`)             // remove multiple spaces by a single
	re_single_q := regexp.MustCompile(`'`) // remove single quote

	for _, line := range lines {
		//line = re_lt_spaces.ReplaceAllString(line, "$1")
		//line = re_mul_spaces.ReplaceAllString(line, "")
		line = re_single_q.ReplaceAllString(line, "")
		if len(line) > cw_dim {
			continue // word  is too long
		}
		line = strings.ToUpper(line)
		//fmt.Println(line)

		wordlist[len(line)] = append(wordlist[len(line)], line)

	}
	fmt.Println("parsed")
}

func File2Lines(filepath string) []string {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
