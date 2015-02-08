// Global constants
package globals

const APPNAME = "HEYGO"
const VERSION = "0.2"
const DATE = "2014-12-14"
const AUTHOR = "Julien CHAUMONT"
const WEBSITE = "https://julienc.io"

const DATABASE = "heygo.db"
const SALT_LENGTH = 15
const ADMIN_GROUP_ID = 1

var CONFIGURATION Configuration
var LoadConfiguration = make(chan bool)
