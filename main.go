package main

import (
	"container/list"
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

func (s *Stack) Push(item rune) int {
	for e := s.Items.Front(); e != nil; e = e.Next() {
		if e.Value.(rune) == item {
			return -1
		}
	}

	s.Items.PushBack(item)
	if s.Items.Len() == s.Length+1 {
		s.Items.Remove(s.Items.Front())
		return 1
	}

	return 0
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

	duplicates := map[rune]int{}
	lengths := []int{}
	stack := Stack{Length: k, Items: list.List{}}
	for _, c := range s {
		switch stack.Push(c) {
		case -1:
			duplicates[c] += 1
		case 0:
			duplicates[c] = 0
		case 1:
			length := stack.Items.Len()
			for e := stack.Items.Front(); e != nil; e = e.Next() {
				length += duplicates[e.Value.(rune)]
			}
			lengths = append(lengths, length)
			break
		}
	}

	if stack.Items.Len() < k {
		fmt.Fprintf(w, "%d", stack.Items.Len())
		return
	}

	m := 0
	for i, e := range lengths {
		if i == 0 || e < m {
			m = e
		}
	}

	fmt.Fprintf(w, "%d", m)
}

// Given an integer k and a string s, find the length of the longest substring that contains
// at most k distinct characters. For example, given s = "abcba" and k = 2, the longest substring
// with k distinct characters is "bcb". Implement in an HTTP endpoint using the standard library.
// Use bluele/gcache to cache evaluated results.
//
// We can keep track of the length of the current substring while we're only seeing two letters
// in a stack, and when we encounter a new one, we pop and start the counter over. We could also use
// slicing and only peek the last two letters that have appeared and keep a big list. But that's not
// very efficient or elegant. The stack would be initialized with k as the length. It would need to be
// a stack that only contains unique values so we don't end up pushing more than one of the same
// character. This could be implemented logically in the push function.
//
// If k > number of distinct characters in the string, we have to rely on the length of Stack.Items.
//
// 01 : read a
//      length = 1
// 02 : read b
//      length = 2
// 03 : read c
//      drop a
//
// substring length = length of registered characters + number of duplicates seen
func main() {
	http.HandleFunc("/", SubstringLengthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
