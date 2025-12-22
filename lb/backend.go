package main

type Backend struct {
	Addr    string
	Healthy bool
	Weight  int
}
