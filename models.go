package main

import "time"

type Task struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	Deadline    time.Time `json:"deadline,omitempty"`
}