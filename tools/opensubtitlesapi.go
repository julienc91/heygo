// Interactions with the OpenSubtitles API
package tools

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/kolo/xmlrpc"
	"heygo/globals"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var username, password, useragent string
var url = "http://api.opensubtitles.org/xml-rpc"
var token = ""
var client *xmlrpc.Client

func init() {
	client, _ = xmlrpc.NewClient(url, nil)
}

func InitFromConfiguration() {
	username = globals.CONFIGURATION.OpensubtitlesLogin
	password = globals.CONFIGURATION.OpensubtitlesPassword
	useragent = globals.CONFIGURATION.OpensubtitlesUseragent
}

// Defer this function's call in the main function
func OpenSubtitlesClose() {
	Logout()
	client.Close()
}

// Compute the hash of a file, and returns it plus its size
// Adapted from https://github.com/oz/osdb
func Hash(filename string) (string, uint64, error) {

	file, err := os.Open(filename)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", 0, err
	}

	var size = uint64(stat.Size())
	var hash = size
	var b = make([]byte, 131072)

	_, err = file.Read(b[:65536])
	if err != nil {
		return "", 0, err
	}

	_, err = file.Seek(-65536, 2)
	if err != nil {
		return "", 0, err
	}

	_, err = file.Read(b[65536:])
	if err != nil {
		return "", 0, err
	}

	for i := 0; i < 16384; i++ {
		hash += binary.LittleEndian.Uint64(b[i*8 : i*8+8])
	}

	return strconv.FormatUint(hash, 16), size, nil
}

// Create a token
func Login() error {

	var result = make(map[string]interface{})
	err := client.Call("LogIn", []interface{}{username, password, "", useragent}, &result)
	if err != nil {
		return err
	}
	if result["status"] == "200 OK" {
		token = result["token"].(string)
	} else {
		return errors.New(result["status"].(string))
	}
	return nil
}

// Delete the current token
func Logout() error {
	return client.Call("LogIn", token, nil)
}

// Search subtitles
func SearchSubtitles(hash string, size uint64) (string, error) {

	Login()
	defer Logout()
	var result = make(map[string]interface{})
	err := client.Call("SearchSubtitles", []interface{}{token, []interface{}{map[string]interface{}{"moviehash": hash, "moviebytesize": size, "sublanguageid": "fre"}}}, &result)
	if err != nil {
		return "", err
	}

	data1, ok := result["data"].([]interface{})
	if !ok || len(data1) == 0 {
		return "", nil
	}

	data, ok := data1[0].(map[string]interface{})
	if !ok {
		return "", nil
	}

	var url = data["ZipDownloadLink"].(string)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var buffer = new(bytes.Buffer)
	n, err := buffer.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	var reader = bytes.NewReader(buffer.Bytes())
	zipReader, err := zip.NewReader(reader, n)

	if err != nil {
		return "", err
	}

	for _, f := range zipReader.File {
		if path.Ext(f.Name) == ".srt" {

			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			var buffer = new(bytes.Buffer)
			_, err = io.Copy(buffer, rc)
			if err != nil {
				return "", err
			}
			var subtitles = srtToVtt(buffer.String())
			return subtitles, nil
		}
	}

	return "", err
}

// Convert SubRip subtitles to WebVtt subtitles
func srtToVtt(srt string) string {

	var srtLines = strings.Split(srt, "\n")
	var r = regexp.MustCompile("\\d+:\\d+:\\d+,\\d+\\s-->\\s\\d+:\\d+:\\d+,\\d+")

	for i := range srtLines {
		if r.MatchString(srtLines[i]) {
			srtLines[i] = strings.Replace(srtLines[i], ",", ".", -1)
		}
	}
	return "WEBVTT\n\n" + strings.Join(srtLines, "\n")
}
