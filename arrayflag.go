package main

import "strings"

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ";")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}
