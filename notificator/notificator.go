package notificator

import (
	"scheduler/utility"
	"scheduler/customtimer"
	"scheduler/domainserver"
	"sync"
	"time"
	"log"
	"os/exec"
)

type Notifier struct{
	Database map[string]dbRecord
	mu sync.RWMutex
}

type dbRecord struct{
	Done chan struct{}
	Message utility.Notification
	timer * customtimer.CTimer
}

func (d * dbRecord)notify(){
	duration,_ := time.ParseDuration(d.Message.Delay)
	var delay = uint32(duration.Seconds())
	if d.Message.DueTime.Unix() - time.Now().Unix() < int64(duration.Seconds()){
		delay = uint32(d.Message.DueTime.Unix() - time.Now().Unix())
	}
	d.timer.Set(delay)
	d.timer.Run()
	for{
		select{
		case d.timer.Expired:
			if d.Message.DueTime.Unix() > time.Now().Unix(){
				exec.Command("notify-send","-i", "/usr/share/icons/gnome/32x32/emotes/face-angel.png", d.Message.Title, d.Message.Body)
				d.timer.Set(delay)
				d.timer.Run()
			}else{
				exec.Command("notify-send", "-i", "/usr/share/icons/gnome/32x32/emotes/face-monkey.png", d.Message.Title, d.Message.Body, "-u", "critical")
				return
			}
		case d.Done:
			return
		}
	}
}

func New()Notifier{
	return Notifier{Database:make(map[string]dbRecord)}
}

func (n * Notifier)Start(){
	server := domainserver.UnixSocketServer{"/tmp/test",n.messageHandler}
	server.OpenSocket()
}

func (n * Notifier)messageHandler(data utility.ChannelData)string{
	var output string
	if data.Err != nil{
		log.Println("Something went wrong, debug program")
		return ""

	}
	if data.Command == "ls"{

		n.mu.Lock()
		defer n.mu.Unlock()
		for _,element := range n.Database{
			output += msgToString(element.Message)
		}
		return output

	}else if data.Command != ""{
		n.mu.Lock()
		defer n.mu.Unlock()
		if _,ok := n.Database[data.Command];ok{
			n.Database[data.Command].Done <- struct{}{}
			delete(n.Database, data.Command)
			return data.Command + " deleted succesfully"
		}else{
			return data.Command + "does not exist"
		}
	}else{
		n.mu.Lock()
		defer n.mu.Unlock()
		if _,ok := n.Database[data.Title];ok{
			return "Task" + data.Title + " already exists, try creating task with another name"
		}else{
			n.Database[data.Title] = dbRecord{make(chan struct{}), data.Notification, customtimer.Init()}
			go n.Database[data.Title].notify()
			return "notification set: " + data.Title
		}
	}
}

func msgToString(n utility.Notification)string{
	return "title: " + n.Title + " body: " + n.Body + " due: "+n.DueDate+" delay: "+ n.Delay+"\n"
}