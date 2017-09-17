package mediator

import (
	"github.com/calendar-scheduler/utility"
	"github.com/calendar-scheduler/customtimer"
	"github.com/calendar-scheduler/domainserver"
	"sync"
	"time"
	"log"
	"os/exec"
	"encoding/json"
)

type Mediator struct{
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
		case <-d.timer.Expired:
			log.Println("notify: ", d.Message.Title)
			if d.Message.DueTime.Unix() > time.Now().Unix(){
				cmd := exec.Command("notify-send","-i", "/usr/share/icons/gnome/32x32/emotes/face-angel.png", d.Message.Title, d.Message.Body)
				err := cmd.Run()
				if err != nil{
					log.Println(err.Error())
				}
				d.timer.Set(delay)
				d.timer.Run()
			}else{
				cmd := exec.Command("notify-send", "-i", "/usr/share/icons/gnome/32x32/emotes/face-monkey.png", d.Message.Title, d.Message.Body, "-u", "critical")
				err := cmd.Run()
				if err != nil{
					log.Println(err.Error())
				}
				return
			}
		case <-d.Done:
			log.Println("expired:", d.Message.Title)
			return
		}
	}
}

func New() Mediator {
	return Mediator{Database:make(map[string]dbRecord)}
}

func (m *Mediator)Start(){
	server := domainserver.UnixSocketServer{"/tmp/test", m.messageHandler}
	server.OpenSocket()
}

func (m *Mediator)messageHandler(data utility.ChannelData)string {
	var output string
	if data.Err != nil {
		log.Println("Something went wrong, debug program\n")
		return ""

	}
	if data.Command == "ls" {

		m.mu.Lock()
		defer m.mu.Unlock()
		for _, element := range m.Database {
			output += msgToString(element.Message)
		}
		return output

	} else if data.Command == "js"{
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, element := range m.Database {	//	TODO add separator
			tmp,_ := json.Marshal(element.Message)
			output += string(tmp)
		}
		return output
	}else if data.Command != ""{
		m.mu.Lock()
		defer m.mu.Unlock()
		if _,ok := m.Database[data.Command];ok{
			m.Database[data.Command].Done <- struct{}{}
			delete(m.Database, data.Command)
			return data.Command + " deleted succesfully\n"
		}else{
			return data.Command + " does not exist\n"
		}
	}else{
		m.mu.Lock()
		defer m.mu.Unlock()
		if _,ok := m.Database[data.Title];ok{
			return "Task" + data.Title + " already exists, try creating task with another name\n"
		}else{
			m.Database[data.Title] = dbRecord{make(chan struct{}), data.Notification, customtimer.Init()}
			i,_ := m.Database[data.Title] //	map returns rvalue, so we should store it to run a pointer method
			go (&i).notify()
			return "notification set: " + data.Title + "\n"
		}
	}
}

func msgToString(n utility.Notification)string{
	return "title: " + n.Title + " body: " + n.Body + " due: "+n.DueDate+" delay: "+ n.Delay+"\n"
}