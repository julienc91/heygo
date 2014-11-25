// Global constants
package globals

const APPNAME = "HEYGO"
const VERSION = "0.1"
const DATE = "2014-10-01"
const AUTHOR = "Julien CHAUMONT"
const WEBSITE = "https://julienc.io"

const DATABASE = "heygo.db"
const SALT_LENGTH = 15
const ADMIN_GROUP_ID = 1

type Configuration struct {
	Domain                 string `json:"domain"`
	Port                   string `json:"port"`
	OpensubtitlesLogin     string `json:"opensubtitles_login"`
	OpensubtitlesPassword  string `json:"opensubtitles_password"`
	OpensubtitlesUseragent string `json:"opensubtitles_useragent"`
}

var CONFIGURATION Configuration
var LoadConfiguration = make(chan bool)
