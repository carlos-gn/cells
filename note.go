package main

import "time"

type Note struct {
	createdAt   time.Time
	title       string
	description string
	content     string
	id          int
}

func (n Note) FilterValue() string {
	return n.title
}

func (n Note) Title() string {
	return n.title
}

func (n Note) Description() string {
	return n.description
}
