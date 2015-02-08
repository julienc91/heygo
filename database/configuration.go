package database

import (
	"github.com/julienc91/heygo/globals"
)

func LoadConfiguration() error {

	var query = "SELECT key, value FROM configuration;"
	stmt, err := db.Preparex(query)
	if err != nil {
		return err
	}

	rows, err := stmt.Queryx()
	if err != nil {
		return err
	}

	for rows.Next() {
		var m = make(map[string]interface{})
		err = rows.MapScan(m)
		if err != nil {
			return err
		}
		switch string(m["key"].([]uint8)) {
		case "domain":
			globals.CONFIGURATION.Domain = string(m["value"].([]uint8))
		case "port":
			globals.CONFIGURATION.Port = string(m["value"].([]uint8))
		case "opensubtitles_login":
			globals.CONFIGURATION.OpensubtitlesLogin = string(m["value"].([]uint8))
		case "opensubtitles_password":
			globals.CONFIGURATION.OpensubtitlesPassword = string(m["value"].([]uint8))
		case "opensubtitles_useragent":
			globals.CONFIGURATION.OpensubtitlesUseragent = string(m["value"].([]uint8))
		}
	}
	return nil
}

func UpdateConfiguration(configuration globals.Configuration) error {

	updateConfigurationRow := func(key string, value interface{}) error {
		var query = "UPDATE configuration SET value=? WHERE key=?;"
		var params = []interface{}{value, key}
		return insertDb(query, params)
	}

	if err := updateConfigurationRow("domain", configuration.Domain); err != nil {
		return err
	}
	if err := updateConfigurationRow("port", configuration.Port); err != nil {
		return err
	}
	if err := updateConfigurationRow("opensubtitles_login", configuration.OpensubtitlesLogin); err != nil {
		return err
	}
	if err := updateConfigurationRow("opensubtitles_password", configuration.OpensubtitlesPassword); err != nil {
		return err
	}
	if err := updateConfigurationRow("opensubtitles_useragent", configuration.OpensubtitlesUseragent); err != nil {
		return err
	}

	globals.LoadConfiguration <- true

	return nil
}
