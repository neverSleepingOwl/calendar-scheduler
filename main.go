package main

import (
	"scheduler/notificator"
	/*"time"
	"github.com/coreos/go-systemd/daemon"
	"log/syslog"
	"log"*/
)

/*func init(){
	logWriter, e := syslog.New(syslog.LOG_NOTICE, "myprog")
	if e == nil {
		log.SetOutput(logWriter)
	}
}*/

func main(){
	/*daemon.SdNotify(false, "READY=1")
	go func() {
		interval, err := daemon.SdWatchdogEnabled(false)
		if err != nil || interval == 0 {
			return
		}
		for {
			daemon.SdNotify(false, "WATCHDOG=1")
			time.Sleep(interval / 3)
		}
	}()*/
	n := notificator.New()
	n.Start()
}