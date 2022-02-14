// Given an integer k and a string s, find the length of the longest substring that contains
// at most k distinct characters. For example, given s = "abcba" and k = 2, the longest substring
// with k distinct characters is "bcb". Implement as an HTTP endpoint using the standard library.
// Use bluele/gcache to cache evaluated results.

package main

import (
	"container/list"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Stack struct {
	Length int
	Items  list.List
}

// Push a rune onto the back of the stack if it does not already exist on it,
// and remove the front value if capacity is exceeded, returning exit codes of 0 and 1, respectively.
// Return the given rune's code as an integer if popped.
//
// Note that this function returns an integer, which could be a rune code
// to be used in the default case, the event of which is a new substring being created.
// We assume that the input string will never contain U+0000 (NULL) or U+0001 (Start of Heading).
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

// Get the length of the longest substring in string s with at most k distinct characters.
func getLongestSubstringLength(k int, s string) int {
	// Add a terminator so we can push the string all the way through.
	s += "\000"

	// Iterate through the input string and map the lengths of encountered substrings.
	// To determine length, we sum the number of characters in the current substring
	// with the number of times we've encountered duplicates.
	//
	// Note that we use default to catch any returned rune code as states 0 and 1 are the
	// only other possible returns from Stack.Push().
	duplicates := map[rune]int{}
	lengths := []int{}
	stack := Stack{Length: k, Items: list.List{}}
	for _, c := range s {
		result := stack.Push(c)
		switch result {
		case 0:
			duplicates[c] += 1
		case 1:
			duplicates[c] = 0
		default:
			length := stack.Items.Len()
			for _, d := range duplicates {
				length += d
			}
			delete(duplicates, rune(result))
			lengths = append(lengths, length)
			break
		}
	}

	if stack.Items.Len() < k {
		return len(s) - 1
	}

	max := 0
	for _, l := range lengths {
		if max < l {
			max = l
		}
	}

	return max
}

func SubstringLengthHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	// Expect the request body to be in the format k,s
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

	// Create a strict-length hash to efficiently support long inputs
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

func a() {
	http.HandleFunc("/", SubstringLengthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
