package model

import "time"

type Estimation struct {
	Id        int
	UserId    uint32
	Segment   string
	CreatedAt time.Time
}
