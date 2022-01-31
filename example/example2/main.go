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
		sticky.Constructor(example.NewInMemoryRepository),
		sticky.Constructor(example.NewAPIRepository),
		sticky.Constructor(func(r *example.InMemoryRepository) *example.UserService {
			return example.NewUserService(r)
		}, sticky.Tag("memory-service")),
		sticky.Constructor(func(r *example.APIRepository) *example.UserService {
			return example.NewUserService(r)
		}, sticky.Tag("api-service")),
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	s, err := sticky.Resolve[*example.UserService](c, sticky.Tag("memory-service"))
	if err != nil {
		panic(err)
	}
	name := s.FindName("001")
	fmt.Printf("id=%v, name=%v\n", "001", name)

	s, err = sticky.Resolve[*example.UserService](c, sticky.Tag("api-service"))
	if err != nil {
		panic(err)
	}
	name = s.FindName("001")
	fmt.Printf("id=%v, name=%v\n", "001", name)
}
