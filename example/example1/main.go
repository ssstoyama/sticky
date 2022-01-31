package main

import (
	"fmt"

	"github.com/ssstoyama/sticky"
	"github.com/ssstoyama/sticky/example"
)

var c sticky.Container

func init() {
	c = sticky.New()

	err := sticky.Register(c,
		sticky.Constructor(example.NewInMemoryRepository, sticky.Implements[example.Repository]()),
		sticky.Constructor(example.NewUserService),
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	s, err := sticky.Resolve[*example.UserService](c)
	if err != nil {
		panic(err)
	}

	name := s.FindName("001")
	fmt.Printf("id=%v, name=%v\n", "001", name)
}
