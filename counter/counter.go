package counter

import (
	"bufio"
	"sort"
	"strings"
)

type WCounter interface {
	Reset()
	AddSentence(string)
	GetStats() []StatElem
}

// CreateCounter returns new object implementing WCounter interface
func CreateCounter() WCounter {
	return &wcounter{stats: map[byte]map[string][]int{}, wcount: 0}
}

type StatElem struct {
	Word      string
	Count     int
	Sentences []int
}

type wcounter struct {
	// to avoid creation of large single map in case of long input text word statistics
	// are spread into multiple maps indexed by first letter of the words
	// i.e. statstics for words "safety" and "surviva" will be stored in one map: stats['s']["safety"] and stats['s']["survival"]
	stats  map[byte]map[string][]int
	scount int
	wcount int
}

// Reset re-initializes internal attributes
func (w *wcounter) Reset() {
	w.stats = map[byte]map[string][]int{}
	w.scount = 0
	w.wcount = 0
}

func (w *wcounter) AddSentence(snt string) {
	if len(snt) == 0 {
		// skip empty sentence
		return
	}
	parts := strings.Split(snt, " ")
	for _, prt := range parts {
		// set letter to lowercase
		currWord := strings.ToLower(prt)
		// trim characters that are not expected to be in special words
		currWord = strings.Trim(currWord, ",;:?!\"\n\r")
		// check if it's a special word and process it only
		if w.processSpecial(currWord) {
			continue
		}
		// trim another non-word characters
		currWord = strings.Trim(currWord, ".'")
		// save word in stats
		w.addWord(currWord)
	}
	// increment number of sentences injected
	w.scount++
}

// GetStats produces slice of
func (w *wcounter) GetStats() []StatElem {
	if w.wcount == 0 {
		return []StatElem{}
	}
	result := make([]StatElem, 0, w.wcount)

	letters := make([]byte, 0, len(w.stats))
	for ltr := range w.stats {
		letters = append(letters, ltr)
	}
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })

	for _, ltr := range letters {
		statPart := w.stats[ltr]
		// sort words for this letter (second stat "index")
		keys := make([]string, 0, len(statPart))
		for k := range statPart {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			stats := statPart[k]
			result = append(result, StatElem{Word: k, Count: len(stats), Sentences: stats})
		}
	}

	return result
}

// addWord function adds statistics for given word
// it is assumed that word is already in correct form (i.e. set to lowercase)
func (w *wcounter) addWord(wrd string) {
	if len(wrd) == 0 {
		// skip empty word
		return
	}
	idx := wrd[0]
	// create new sub-map if needed
	if _, ok := w.stats[idx]; !ok {
		w.stats[idx] = map[string][]int{}
	}
	// save the word and update counter if needed
	wStat, ok := w.stats[idx][wrd]
	if !ok {
		// create new slice with occurences
		w.stats[idx][wrd] = []int{(w.scount + 1)}
		// update word counter
		w.wcount++
	} else {
		// attach new sentence index at the end of existing slice
		w.stats[idx][wrd] = append(wStat, (w.scount + 1))
	}
}

// processSpecial process given word in special way if it's added to 'specialWords' and returns true
// If given word is "not special" function returns false without any statistics modification
func (w *wcounter) processSpecial(wrd string) bool {
	newContent, ok := specialWords[wrd]
	if !ok {
		return false
	}
	for _, newWord := range newContent {
		w.addWord(newWord)
	}
	return true
}

// specialWords contains words ar abbreviations that should be treated in special way:
// - added to stats with dots like "mrs." or "mr."
// - added to stats as mutliple words after expansion like "i am" instead of "i'm" and "you are" instead of you're
var specialWords = map[string][]string{
	"i'm":     {"i", "am"},
	"you're":  {"you", "are"},
	"he's":    {"he", "is"},
	"she's":   {"she", "is"},
	"it's":    {"is", "is"},
	"we're":   {"we", "are"},
	"they're": {"they", "are"},
	"mrs.":    {"mrs."},
	"mr.":     {"mr."},
	"i.e.":    {"i.e."},
	"etc.":    {"etc."},
}

func SentenceSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	bufLen := len(data)

	if bufLen == 0 {
		// should not happen
		if atEOF {
			return 0, nil, bufio.ErrFinalToken
		}
		return 0, nil, bufio.ErrBadReadCount
	}

	var currIdx int
	lastSpace := 0

	for currIdx = 0; currIdx < bufLen; currIdx++ {
		switch data[currIdx] {
		case '!', '?', ';', '\n':
			return currIdx + 1, data[:currIdx+1], nil
		case ' ':
			lastSpace = currIdx
		case '.':
			if currIdx+1 == bufLen && atEOF {
				// it's a finishing dot
				return currIdx + 1, data, bufio.ErrFinalToken
			}
			if currIdx+1 == bufLen && !atEOF {
				// more characters needed
				return currIdx + 1, nil, nil
			}
			if data[currIdx+1] == '\n' {
				// it's probably an end of sentence
				return currIdx + 2, data[:currIdx+2], nil
			}
			if data[currIdx+1] == ' ' {
				currWord := strings.Trim(string(data[lastSpace:currIdx+1]), " ")
				if _, ok := specialWords[strings.ToLower(currWord)]; ok {
					// it's one of special words with dot
					continue
				}
				return currIdx + 1, data[:currIdx+1], nil
			}
		}
	}

	// not returned till now - check if it's not the final chunk:
	if atEOF {
		return bufLen, data, bufio.ErrFinalToken
	}

	return currIdx + 1, nil, nil
}
