package counter

import (
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
		// check if it's a special word
		// trim all non-word characters
		currWord = strings.Trim(currWord, ".,;:?!'")
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
