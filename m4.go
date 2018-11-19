package main

type M4HttpAction string

const (
	M4Reset M4HttpAction = "config.cgi?mres"
	CpcReset M4HttpAction = "config.cgi?cres"
	Start M4HttpAction = "config.cgi?cctr"
	Run   M4HttpAction = "config.cgi?run"
	Pause M4HttpAction = "config.cgi?chlt"
	Upload M4HttpAction  = "upfile"
	Download M4HttpAction = "sd/"
)


type M4Client struct{
	Action M4HttpAction
	IpClient string
	CpcFilePath string
} 