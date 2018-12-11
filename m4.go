package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// M4HttpAction is struct for url complement according to the action
type M4HttpAction string

var userAgent = "cpcxfer"

// M4 Wifi card http possibles actions
const (
	M4Reset  M4HttpAction = "config.cgi?mres"
	CpcReset M4HttpAction = "config.cgi?cres"
	Start    M4HttpAction = "config.cgi?cctr"
	Mkdir    M4HttpAction = "config.cgi?mdkir="
	Ls       M4HttpAction = "config?ls="
	Cd       M4HttpAction = "config?cd="
	Rm       M4HttpAction = "config?rm="
	Execute  M4HttpAction = "config.cgi?run2="
	Run      M4HttpAction = "config.cgi?run"
	Pause    M4HttpAction = "config.cgi?chlt"
	Upload   M4HttpAction = "upload.html"
	Download M4HttpAction = "sd/"
	Rom      M4HttpAction = "roms.shtml"
)

// M4Client M4 http client with action, address ip client
// and Cpc file path
type M4Client struct {
	Action            M4HttpAction
	IPClient          string
	CpcLocalFilePath  string
	CpcRemoteFilePath string
}

func (m *M4Client) Url() string {
	return "http://" + m.IPClient + "/" + string(m.Action)
}

func PerformHttpAction(req *http.Request) error {
	client := &http.Client{}
	req.Header.Add("user-agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Response from cpc http server differs from 200")
	}
	return nil
}

func (m *M4Client) PauseCpc() error {
	m.Action = Pause
	req, err := http.NewRequest("GET", m.Url(), nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) ResetM4() error {
	m.Action = M4Reset
	req, err := http.NewRequest("GET", m.Url(), nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) ResetCpc() error {
	m.Action = CpcReset
	req, err := http.NewRequest("GET", m.Url(), nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) Download(remotePath string) error {
	m.Action = Download
	fh, err := os.Create(path.Base(remotePath))
	if err != nil {
		return err
	}
	defer fh.Close()
	req, err := http.NewRequest("GET", m.Url()+remotePath, nil)
	req.Header.Add("user-agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Not http status ok ")
	}
	_, err = io.Copy(fh, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (m *M4Client) UploadDirectoryContent(remotePath, localDirectoryPath string) error {
	files, err := ioutil.ReadDir(localDirectoryPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			m.Upload(remotePath, localDirectoryPath+string(filepath.Separator)+file.Name())
		}
	}
	return nil
}

func (m *M4Client) Upload(remotePath, localPath string) error {
	m.Action = Upload
	fh, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer fh.Close()
	if _, err := NewCpcHeader(fh); err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upfile", remotePath+"/"+path.Base(localPath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, fh)
	if err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", m.Url(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Expires", "0")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Http response differs from 200")
	}
	return nil
}

func (m *M4Client) Execute(cpcfile string) error {
	m.Action = Execute
	req, err := http.NewRequest("GET", m.Url()+cpcfile, nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) MakeDirectory(remoteFolder string) error {
	m.Action = Mkdir
	req, err := http.NewRequest("GET", m.Url()+remoteFolder, nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) ChangeDirectory(remoteFolder string) error {
	m.Action = Cd
	req, err := http.NewRequest("GET", m.Url()+remoteFolder, nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) DeleteRom(romNumber int) error {
	m.Action = Rom
	req, err := http.NewRequest("GET", m.Url()+"?rmsl="+strconv.Itoa(romNumber), nil)
	if err != nil {
		return err
	}
	return PerformHttpAction(req)
}

func (m *M4Client) UploadRom(romFilpath, romName string, romId int) error {
	if romId < 0 || romId >= 32 {
		return errors.New("Rom id is not compliant.")
	}
	m.Action = Rom

	fh, err := os.Open(romFilpath)
	if err != nil {
		return err
	}
	defer fh.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadedfile", "rom.bin")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, fh)
	if err != nil {
		return err
	}
	slotNumW, err := writer.CreateFormField("slotnum")
	if err != nil {
		return err
	}
	slotNumW.Write([]byte(fmt.Sprintf("%d", romId)))

	slotNameW, err := writer.CreateFormField("slotname")
	if err != nil {
		return err
	}
	slotNameW.Write([]byte(romName))

	writer.Close()

	req, err := http.NewRequest("POST", m.Url(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Expires", "0")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Http response differs from 200")
	}
	return nil
}

func (m *M4Client) Ls(remotePath string) (string, error) {
	m.Action = Ls
	client := &http.Client{}
	req, err := http.NewRequest("GET", m.Url(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("user-agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Response from cpc http server differs from 200")
	}
	return string(body), nil
}
