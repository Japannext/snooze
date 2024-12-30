package errors

// Stop the pipeline for one alert.
type StopPipeline struct{}

func (err *StopPipeline) Error() string {
	return ""
}
