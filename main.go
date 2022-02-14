package main

import "github.com/bluele/gcache"

var Cache gcache.Cache

func init() {
	Cache = gcache.New(10).Build()
}

func main() {
	// A()
	// B()
	C()
}
