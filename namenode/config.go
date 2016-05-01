// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package namenode

import (
	"strings"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/c4pt0r/cfg"
)

type Config struct {
	addr        string
	logFile     string
	logLevel    string
	logFileSize uint64
}

func LoadConf(configFile string) (*Config, error) {
	c := cfg.NewCfg(configFile)
	if err := c.Load(); err != nil {
		log.PanicErrorf(err, "load config '%s' failed", configFile)
	}

	conf := &Config{}
	conf.addr, _ = c.ReadString("addr", "172.0.0.1:9000")
	if len(conf.addr) == 0 {
		log.Panicf("invalid config: addr is missing in %s", configFile)
	}
	conf.logFile, _ = c.ReadString("log_file", "")
	if len(conf.logFile) == 0 {
		log.Panicf("invalid config: log file is missing in %s", configFile)
	}
	conf.logLevel, _ = c.ReadString("log_level", "")
	if len(conf.logLevel) == 0 {
		log.Panicf("invalid config: need zk entry is missing in %s", configFile)
	}
	loadConfInt := func(entry string, defval uint64) uint64 {
		v, _ := c.ReadUInt64(entry, defval)
		if v < 0 {
			log.Panicf("invalid config: read %s = %d", entry, v)
		}
		return v
	}

	conf.logFileSize = loadConfInt("log_file_size", 1024*1024*1024)

	return conf, nil
}
