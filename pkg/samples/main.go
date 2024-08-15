package samples

func Run() error {
	if err := runSyslogSamples(); err != nil {
		return err
	}
	return nil
}
