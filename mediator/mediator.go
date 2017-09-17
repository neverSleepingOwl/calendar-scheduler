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
	//Database, which stores notification entities, which has it's own channel to stop sending
	//notification and own delay timer.
	Database map[string]dbRecord
	//Mutex to protect database
	mu sync.RWMutex
}


//Structure, representing single activity to notify
type dbRecord struct{
	//Chanel for force-kill notification
	Done chan struct{}
	//Notification message itself
	Message utility.Notification
	//Timer to notify after a given delay
	timer * customtimer.CTimer
}

func (d * dbRecord)notify(){
	duration,_ := time.ParseDuration(d.Message.Delay)

	var delay = uint32(duration.Seconds())

	if d.Message.DueTime.Unix() - time.Now().Unix() < int64(duration.Seconds()){
		delay = uint32(d.Message.DueTime.Unix() - time.Now().Unix())	//	if delay is more then time remaining to task expires time then notify once
	}

	d.timer.Set(delay)	//	set and run timer
	d.timer.Run()

	for{
		select{
		case <-d.timer.Expired:	//	repeat notification each n <time interval> means:
			// wait for timer, and reset timer with delay.
			log.Println("notify: ", d.Message.Title)
			//some logging for debug purposes
			if d.Message.DueTime.Unix() > time.Now().Unix(){
				//If time hasn't expired notify normally each <time interval>
				cmd := exec.Command("notify-send","-i", "/usr/share/icons/gnome/32x32/emotes/face-angel.png", d.Message.Title, d.Message.Body)
				err := cmd.Run()	//	set and run command

				if err != nil{
					log.Println(err.Error())	//	log if something bad happens
				}

				d.timer.Set(delay)	//	reset and run timer
				d.timer.Run()

			}else{
				//if task has expired send one critical notification and then stop notifying
				cmd := exec.Command("notify-send", "-i", "/usr/share/icons/gnome/32x32/emotes/face-monkey.png", d.Message.Title, d.Message.Body, "-u", "critical")
				err := cmd.Run()
				if err != nil{
					log.Println(err.Error())
				}
				return
			}
		case <-d.Done:
			//if any other goroutine needs to stop this notifier exit function
			log.Println("expired:", d.Message.Title)
			return
		}
	}
}


//New - constructor for core class
func New() Mediator {
	return Mediator{Database:make(map[string]dbRecord)}
}


//Start - start listening to unix socket and send notifications
func (m *Mediator)Start(){
	server := domainserver.UnixSocketServer{"/tmp/test", m.messageHandler}
	server.OpenSocket()
}


//messageHandler is called in unix domain socket server on each message received
//starts new notifiers, handles commands and so on
func (m *Mediator)messageHandler(data utility.ChannelData)string {
	var output string

	if data.Err != nil {	//	errors are being handled in server part
		log.Println("Something went wrong, debug program\n")
		return ""

	}
	if data.Command == "ls" {	//	command, list all tasks in human-readable view

		m.mu.Lock()		//	protect database with mutex to avoid race or unpredicted behavior
		defer m.mu.Unlock()

		for _, element := range m.Database {	//	iterate through all map
			//	run utility function, which converts gnome notification to human-readable string
			output += msgToString(element.Message)
		}

		return output

	} else if data.Command == "js"{	//	command, send all tasks in json format

		m.mu.Lock()
		defer m.mu.Unlock()

		for _, element := range m.Database {
			tmp,_ := json.Marshal(element.Message)
			output += string(tmp)

		}

		return output

	}else if data.Command != ""{	//	delete command, WARNING:you can't create tasks with titles ls or js
		//	if command isn't listed higher, delete record with a given title
		m.mu.Lock()
		defer m.mu.Unlock()

		if _,ok := m.Database[data.Command];ok{	//	if record with the following title exists
			m.Database[data.Command].Done <- struct{}{}	//	close notifier goroutine
			delete(m.Database, data.Command)	//	delete task
			return data.Command + " deleted succesfully\n"	//	send help message
		}else{
			return data.Command + " does not exist\n"
		}
	}else{	//	if message received isn't a command, create new record
		m.mu.Lock()
		defer m.mu.Unlock()

		if _,ok := m.Database[data.Title];ok{	//	avoid overwriting tasks
			return "Task" + data.Title + " already exists, try creating task with another name\n"
		}else{
			m.Database[data.Title] = dbRecord{make(chan struct{}), data.Notification, customtimer.Init()}

			i,_ := m.Database[data.Title] //	map returns rvalue, so we should store it to run a pointer method
			go (&i).notify()	//	note that in go, function, changing class value must be called like here
			//call notificator loop, which sends notifications each <period of time>
			return "notification set: " + data.Title + "\n"
		}
	}
}


//msgToString - Utility Notification to string function
func msgToString(n utility.Notification)string{
	return "title: " + n.Title + " body: " + n.Body + " due: "+n.DueDate+" delay: "+ n.Delay+"\n"
}