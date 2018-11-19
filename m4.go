package main

// M4HttpAction is struct for url complement according to the action
type M4HttpAction string

// M4 Wifi card http possibles actions
const (
	M4Reset  M4HttpAction = "config.cgi?mres"
	CpcReset M4HttpAction = "config.cgi?cres"
	Start    M4HttpAction = "config.cgi?cctr"
	Run      M4HttpAction = "config.cgi?run"
	Pause    M4HttpAction = "config.cgi?chlt"
	Upload   M4HttpAction = "upfile"
	Download M4HttpAction = "sd/"
)

// M4Client M4 http client with action, address ip client
// and Cpc file path
type M4Client struct {
	Action      M4HttpAction
	IPClient    string
	CpcFilePath string
}
