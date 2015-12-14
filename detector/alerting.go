package detector

// Start several goroutines to wait for detected metrics, then check each
// metric with all the rules, the configured shell command will be executed
// once a rule is hit.
func (d *Detector) startAlertingWorkers() {
	for i := 0; i < d.cfg.Alerter.Workers; i++ {
		go d.alertingWork()
	}
}

func (d *Detector) alertingWork() {
	_ = <-d.rc
	// FIXME
}
