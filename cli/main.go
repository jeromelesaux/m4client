package main

import (
	"flag"
	"fmt"
	"github.com/jeromelesaux/m4client/m4"
	"os"
)

var (
	host       = flag.String("host", "", "Ip V4 or hostname of the your CPC")
	infile     = flag.String("file", "", "file path of the file to get or to send")
	remotePath = flag.String("remotepath", "", "remote path where to get or send your file")
	ls         = flag.Bool("ls", false, "get the current remote path on your m4.")
	cd         = flag.String("cd", "", "change the current remote path on your m4.")
)

func main() {
	flag.Parse()

	if *host == "" {
		fmt.Fprintf(os.Stderr, "Cannot contact M4 without its hostname or IP\n")
		os.Exit(-1)
	}

	if *ls {
		client := m4.M4Client{
			Action:   m4.Ls,
			IPClient: *host,
		}
		rpath, err := client.Ls()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while getting remote path with host (%s) error :%v\n", *host, err)
		} else {
			fmt.Fprintf(os.Stdout, "Remote path (%s) host (%s)\n", rpath, *host)
		}
	}
	if *cd != "" {
		client := m4.M4Client{
			Action:            m4.Cd,
			IPClient:          *host,
			CpcRemoteFilePath: *cd,
		}
		if err := client.ChangeDirectory(); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot change the directory on the M4 (%s) error :%v\n", *host, err)
		}
	}
}
