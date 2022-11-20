package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wcounter_addWord(t *testing.T) {
	testCounter := wcounter{stats: map[byte]map[string][]int{}, wcount: 0}

	testCounter.addWord("test")
	assert.Equal(t, 1, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 0, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 1, len(testCounter.stats), "Invalid number of base indices")

	testCounter.addWord("tester")
	assert.Equal(t, 2, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 0, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 1, len(testCounter.stats), "Invalid number of base indices")

	testCounter.addWord("test")
	assert.Equal(t, 2, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 0, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 1, len(testCounter.stats), "Invalid number of base indices")

	testCounter.addWord("assignment")
	assert.Equal(t, 3, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 0, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 2, len(testCounter.stats), "Invalid number of base indices")

	testCounter.addWord("interview")
	assert.Equal(t, 4, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 0, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 3, len(testCounter.stats), "Invalid number of base indices")
}

func Test_wcounter_Reset(t *testing.T) {
	testCounter := wcounter{stats: map[byte]map[string][]int{}, wcount: 0}

	testCounter.AddSentence("I am tester")

	assert.Equal(t, 3, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 1, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 3, len(testCounter.stats), "Invalid number of base indices")

	testCounter.Reset()
	assert.Equal(t, 0, testCounter.wcount, "Invalid word counter")
	assert.Equal(t, 0, testCounter.scount, "Invalid sentence counter")
	assert.Equal(t, 0, len(testCounter.stats), "Invalid number of base indices")

}

func Test_wcounter_GetStats(t *testing.T) {

	testCounter := CreateCounter()

	stats := testCounter.GetStats()
	assert.Equal(t, 0, len(stats), "Statistics for new counter should be empty")

	testCounter.AddSentence("I am a tester")
	stats = testCounter.GetStats()
	assert.Equal(t, 4, len(stats), "Statistics for new counter should be empty")
	expected := []StatElem{
		{Word: "a", Count: 1, Sentences: []int{1}},
		{Word: "am", Count: 1, Sentences: []int{1}},
		{Word: "i", Count: 1, Sentences: []int{1}},
		{Word: "tester", Count: 1, Sentences: []int{1}},
	}
	assert.Equal(t, expected, stats)

	testCounter.AddSentence("I am not a tester")
	stats = testCounter.GetStats()
	assert.Equal(t, 5, len(stats), "Statistics for new counter should be empty")
}
