// Some useful functions
package tools

import (
	crand "crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienc91/heygo/globals"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// Check if e is in a
func InArray(a []string, e string) bool {
	for _, x := range a {
		if x == e {
			return true
		}
	}
	return false
}

// Check if a file exists
func CheckFilePath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Get all files with the given extension in the given directory
func GetFilesFromSubfolder(subfolder string, extension string, recursive bool) ([]string, error) {

	files, err := ioutil.ReadDir(subfolder)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range files {
		if f.IsDir() && recursive {
			subres, err := GetFilesFromSubfolder(path.Join(subfolder, f.Name()), extension, recursive)
			if err != nil {
				return nil, err
			}
			res = append(res, subres...)
		} else if !f.IsDir() && path.Ext(f.Name()) == extension {
			res = append(res, path.Join(subfolder, f.Name()))
		}
	}
	return res, nil
}

// Get a slug from the given filename
func SlugFromFilename(filename string) string {

	var slug = strings.ToLower(filename)
	var from = "ãàáäâ@ẽèéëêìíïîõòóöôùúüûñç"
	var to = "aaaaaaeeeeeiiiiooooouuuunc"
	for i := 0; i < len(from) && i < len(to); i++ {
		slug = strings.Replace(slug, string(from[i]), string(to[i]), -1)
	}
	var re = regexp.MustCompile("[^\\w]+")
	return re.ReplaceAllString(slug, "_")
}

// Hash function
func Hash(password, salt string) string {
	var hash = sha512.Sum512([]byte(salt + password))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Salt generating function
func SaltGenerator() string {

	mrand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, globals.SALT_LENGTH)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

// UUID generating function
// From: http://play.golang.org/p/4FkNSiUDMg
func UuidGenerator() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// Answer the user's query with the given codes and with a json object constructed
// from 'ret'.
func WriteJsonResult(ret map[string]interface{}, w http.ResponseWriter, code int) {

	val, err := json.Marshal(ret)
	if err != nil {
		ret["err"] = err.Error()
		ret["ok"] = false
		code = http.StatusInternalServerError
	}

	WriteResponse(string(val), w, "application/json", code)
}

// Set http error code and content-type and write response
func WriteResponse(content interface{}, w http.ResponseWriter, contentType string, code int) {

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	fmt.Fprint(w, content)
}
