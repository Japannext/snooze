package decision

// The decision to take when done processing the log.
// Valid Cases:
// * OK: Continue processing
// * Done: The processing did as expected, but need to stop here (silence/drop/snooze/etc)
// Error handling:
// * Retry: Requeue (Nack) and retry later. One database is down for instance.
// * Abort: The process didn't work as expected. Something in the log data is making it fail.
type Decision struct {
	Kind  int
	Stop  bool
	Error error
}

const (
	KindOK int = iota
	KindDone
	KindRetry
	KindAbort
)

func OK() *Decision {
	return &Decision{Kind: KindOK}
}

func Done() *Decision {
	return &Decision{Kind: KindDone, Stop: true}
}

func Retry(err error) *Decision {
	return &Decision{Kind: KindRetry, Error: err, Stop: true}
}

func Abort(err error) *Decision {
	return &Decision{Kind: KindAbort, Error: err, Stop: true}
}
