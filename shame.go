package main

import "time"

type Shame struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time

	URL      string `sql:"size:511"`
	Reponame string
	Path     string
	Author   string
	Content  []byte
}
