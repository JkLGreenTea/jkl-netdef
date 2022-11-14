package interfacies

// ControllerTechnicalWorks - контроллер тех. работ.
type ControllerTechnicalWorks interface {
	Check(host, userAgent string) int64
	StopGlobal()
	StartGlobal(startTime, endTime int64)
	IssueGlobalAccess(host, userAgent string)

	CheckByPage(url, host, userAgent string) int64
	StopByPage(url string)
	StartByPage(startTime, endTime int64, url string)
	IssueByPagesAccess(host, userAgent, url string)
}
