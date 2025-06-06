package models

import "time"

type Problem struct {
	UUID            string    `json:"uuid" db:"uuid" table:"problems"`
	ParentProblemID string    `json:"parentProblemID,omitempty"`
	Name            string    `json:"name" db:"name"`
	Kind            string    `json:"kind"`
	Status          string    `json:"status"`
	StartTime       time.Time `json:"startTime"`
	SolvedTime      time.Time `json:"solvedTime"`
}
