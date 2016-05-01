package main

import (
	"github.com/docopt/docopt-go"
	"github.com/sjarvie/godfs/utils/log"
)

var (
	configFile = "config.ini"
)

var usage = `usage: namenode [-c <config_file>]

options:
   -c	set config file
`

func main() {
	args, err := docopt.Parse(usage, nil, true, "DFS 0.1", true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set config file
	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}
	conf, err := namenode.LoadConf(configFile)
	if err != nil {
		log.PanicErrorf(err, "load config failed")
	}

	checkUlimit(1024)
	runtime.GOMAXPROCS(8)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)

	s := namenode.New(conf)
	defer s.Close()
	go func() {
		<-c
		log.Info("ctrl-c or SIGTERM found, bye bye...")
		s.Close()
	}()

	time.Sleep(time.Second)

	s.Join()
	log.Infof("namenode exit!! :(")
}

func checkUlimit(min int) {
	ulimitN, err := exec.Command("/bin/sh", "-c", "ulimit -n").Output()
	if err != nil {
		log.WarnErrorf(err, "get ulimit failed")
	}

	n, err := strconv.Atoi(strings.TrimSpace(string(ulimitN)))
	if err != nil || n < min {
		log.Panicf("ulimit too small: %d, should be at least %d", n, min)
	}
}
