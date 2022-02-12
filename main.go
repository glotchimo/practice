// Given an integer k and a string s, find the length of the longest substring that contains
// at most k distinct characters. For example, given s = "abcba" and k = 2, the longest substring
// with k distinct characters is "bcb". Implement in an HTTP endpoint using the standard library.
// Use bluele/gcache to cache evaluated results.

package main

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bluele/gcache"
)

var Cache gcache.Cache

type Stack struct {
	Length int
	Items  list.List
}

// Push a rune onto the back of the stack if it does not already exist on it,
// and remove the front value if capacity is exceeded, returning exit codes of 0 and 1, respectively.
// Return the rune's value as an integer if popped.
func (s *Stack) Push(r rune) int {
	for e := s.Items.Front(); e != nil; e = e.Next() {
		if e.Value.(rune) == r {
			return 0
		}
	}

	s.Items.PushBack(r)
	if s.Items.Len() == s.Length+1 {
		p := s.Items.Front()
		s.Items.Remove(p)
		return int(p.Value.(rune))
	}

	return 1
}

func getLongestSubstringLength(k int, s string) int {
	s += "\000" // We need to add a null byte to flush the buffer

	duplicates := map[rune]int{}
	lengths := []int{}
	stack := Stack{Length: k, Items: list.List{}}
	for _, c := range s {
		// Note that the Stack.Push function returns an integer, which could be a rune code
		// to be used in the default case, the event of which is a new substring created.
		// We assume that the input string will never contain U+0000 (NULL) or U+0001 (Start of Heading).
		result := stack.Push(c)
		switch result {
		case 0: // Duplicate found; increment corresponding counter
			duplicates[c] += 1
		case 1: // No change; register character in duplicates map
			duplicates[c] = 0
		default: // New substring started; delete outgoing character entry, register substring length
			length := stack.Items.Len()
			for _, v := range duplicates {
				length += v
			}
			delete(duplicates, rune(result))
			lengths = append(lengths, length)
			break
		}
	}

	if stack.Items.Len() < k {
		return len(s) - 1
	}

	m := 0
	for _, e := range lengths {
		if m < e {
			m = e
		}
	}

	return m
}

func SubstringLengthHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	contents := strings.Split(string(data), ",")
	if len(contents) != 2 {
		w.WriteHeader(400)
	}

	k, err := strconv.Atoi(contents[0])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	s := contents[1]

	hasher := sha256.New()
	hasher.Write([]byte(s))
	hash := hasher.Sum(nil)
	strHash := string(hash)

	length, errGet := Cache.Get(strHash)
	if errGet != nil {
		length = getLongestSubstringLength(k, s)
		if errSet := Cache.Set(strHash, length); errSet != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Cache write failed: %s", err.Error())
			return
		}
	}
	fmt.Fprintf(w, "%d", length)
}

func init() {
	Cache = gcache.New(10).Build()
}

func main() {
	http.HandleFunc("/", SubstringLengthHandler)
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()

	letters := []string{"a", "b", "c", "d", "e", "f"}
	builder := strings.Builder{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000000; i++ {
		builder.WriteString(letters[rand.Intn(len(letters))])
	}
	s := builder.String()

	b := bytes.Buffer{}
	b.WriteString("3," + s)
	req1, _ := http.NewRequest("POST", "http://localhost:8080/", &b)
	res1, err := http.DefaultClient.Do(req1)
	if err != nil {
		panic(err)
	}
	data1, err := io.ReadAll(res1.Body)
	if err != nil {
		panic(err)
	}
	defer res1.Body.Close()
	fmt.Println(string(data1))

	b.WriteString("3," + s)
	req2, _ := http.NewRequest("POST", "http://localhost:8080/", &b)
	res2, err := http.DefaultClient.Do(req2)
	if err != nil {
		panic(err)
	}
	data2, err := io.ReadAll(res2.Body)
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()
	fmt.Println(string(data2))
}
