package models

type TimeRange struct {
	// Start time in milliseconds
	Start uint64 `form:"start"`
	// End time in milliseconds
	End uint64 `form:"end"`
}
