// Interactions with the OpenSubtitles API
package tools

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
func OpensubtitlesHash(filename string) (string, uint64, error) {

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
func SearchSubtitles(hash string, size uint64, imdbId string, lang string) ([]string, error) {

	search := func(token string, params map[string]interface{}) ([]map[string]interface{}, error) {

		var result = make(map[string]interface{})
		err := client.Call("SearchSubtitles", []interface{}{token, []interface{}{params}}, &result)
		if err != nil {
			return nil, err
		}

		list, ok := result["data"].([]interface{})
		if !ok {
			if _, ok := result["data"].(bool); ok {
				return nil, nil
			}
			fmt.Println(result)
			return nil, errors.New("Type assertion error for data")
		}

		var res []map[string]interface{}
		for i := range list {
			data, ok := list[i].(map[string]interface{})
			if ok {
				res = append(res, data)
			}
		}
		if len(res) == 0 {
			return nil, nil
		}

		return res, nil
	}

	downloadArchive := func(url string) (*bytes.Buffer, int64, error) {

		resp, err := http.Get(url)
		if err != nil {
			return nil, 0, err
		}
		defer resp.Body.Close()

		var buffer = new(bytes.Buffer)
		n, err := buffer.ReadFrom(resp.Body)
		if err != nil {
			return nil, 0, err
		}
		return buffer, n, nil
	}

	handleZipFile := func(buffer *bytes.Buffer, n int64) (string, error) {

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

				return TryToUtf8(buffer.String(), lang), nil
			}
		}
		return "", errors.New("No srt subtitles found")
	}

	Login()
	defer Logout()

	subtitlesList, err := search(token, map[string]interface{}{"moviehash": hash, "moviebytesize": size, "sublanguageid": lang})
	if err != nil {
		return nil, err
	} else if subtitlesList == nil && imdbId != "" {
		var imdbNumber int

		if len(imdbId) <= 2 || !strings.HasPrefix(imdbId, "tt") {
			return nil, errors.New("Invalid imdb id")
		} else if imdbNumber, err = strconv.Atoi(imdbId[2:]); err != nil {
			return nil, err
		}

		subtitlesList, err = search(token, map[string]interface{}{"imdbid": imdbNumber, "sublanguageid": lang})
		if err != nil {
			return nil, err
		}
	}

	if subtitlesList == nil {
		return nil, nil
	}

	var res []string
	for _, data := range subtitlesList {

		var url = data["ZipDownloadLink"].(string)
		buffer, n, err := downloadArchive(url)
		if err != nil {
			continue
		}

		subtitles, err := handleZipFile(buffer, n)
		if err != nil {
			continue
		}
		res = append(res, srtToVtt(subtitles))
	}

	return res, nil
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
