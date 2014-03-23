// config
package main

import (
	"flag"
	"fmt"
	"github.com/prezi/go-ini"
	"log"
	"strconv"
)

var config *ini.File

func cfgGetInt(config *ini.File, section, key string, defval int64) (val int64) {
	strval := config.GetWithDefault(section, key, fmt.Sprintf("%d", defval))
	val, err := strconv.ParseInt(strval, 0, 64)
	if err != nil {
		val = defval
	}
	return
}

func ReadConfig(defaultConfig string) *ini.File {
	var configfile = flag.String("config", defaultConfig, "Config file path")
	flag.Parse()
	config, err := ini.LoadFile(*configfile)
	if err != nil {
		config = make(ini.File)
		log.Println("Config data created")
	} else {
		log.Printf("Config data loaded: %+v", config)
	}
	return &config
}
