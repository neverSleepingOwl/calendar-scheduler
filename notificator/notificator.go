package notificator

import (
	"scheduler/utility"
	"scheduler/customtimer"
	"scheduler/domainserver"
	"sync"
	"time"
	"log"
)

type Notifier struct{
	Database map[string]dbRecord
	mu sync.RWMutex
}

type dbRecord struct{
	Done chan struct{}
	Message utility.Notification
	timer customtimer.CTimer
}

func createRecord(){

}

func (d * dbRecord)notify(){

}

func New()Notifier{
	return Notifier{Database:make(map[string]dbRecord)}
}

func (n * Notifier)Start(){
	server := domainserver.UnixSocketServer{"/tmp/test",n.messageHandler}
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
			output += msgToString(element)
		}
		return output

	}else if data.Command != ""{
		n.mu.Lock()
		defer n.mu.Unlock()
		if i,ok := n.Database[data.Command];ok{
			n.Database[data.Command].Done <- struct{}{}
			delete(n.Database, data.Command)
			return data.Command + " deleted succesfully"
		}else{
			return data.Command + "does not exist"
		}
	}else{
		n.mu.Lock()
		defer n.mu.Unlock()
		if i,ok := n.Database[data.Command];ok{
		}
	}
}

func msgToString(n utility.Notification)string{
	return "title: " + n.Title + " body: " + n.Body + " due: "+n.DueDate+" delay: "+ n.Delay+"\n"
}