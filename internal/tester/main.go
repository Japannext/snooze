package tester

func Run() error {
	if err := runSyslogSamples(); err != nil {
		return err
	}
	return nil
}
