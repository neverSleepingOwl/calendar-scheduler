package utility

import (
	"regexp"
	"errors"
	"encoding/json"
	"time"
	"fmt"
)

//package, containing small functions for different purposes:
//parsing configs, parsing messages, checking time for correctness

//{"title":"tag","body":"testing arp requests","due":"14 Sep 17 17:00 GMT","delay":"0.05h"}

type Notification struct{
	Title     string	`json:"title"`     //	notification title
	Body      string	`json:"body"`      //	notification body
	DueDate   string	`json:"due"`       //	due date
	Delay     string    `json:"delay"`	   //  delay before two notifications
	DueTime   time.Time                    //	due date, but has time format
}

//struct for channel output
type ChannelData struct{
	Notification
	Command string
	Err error
}

//NweNotification parse message
func NewNotification(input []byte)(ChannelData) {
	var c ChannelData = ChannelData{Notification:Notification{},Err:nil}
	if string(input)[len(string(input))-1] == '\n'{	//	TODO fix terminating character
		input = input[:len(input)-1]
	}
	command := string(input)

	deleteChecker := regexp.MustCompile(`^(([rR]emove)|([Dd]elete))\s+\w+$`)
	deleteCommand := regexp.MustCompile(`^(([rR]emove)|([Dd]elete))`)
	delimeter:= regexp.MustCompile(`\s`)

	listChecker := "ls"

	if deleteChecker.MatchString(string(input)){
		s := delimeter.ReplaceAllString(command,"")
		s = deleteCommand.ReplaceAllString(s,"")
		return ChannelData{Command:s}
	}else if command == listChecker{
		return ChannelData{Command:"ls"}
	}

	if err:=json.Unmarshal(input,&c.Notification);err == nil{
		c.DueTime,err = time.Parse(time.RFC822,c.DueDate)
		fmt.Println(time.Now().Unix())
		fmt.Println(c.DueTime.Unix())
		if c.DueTime.Unix() < time.Now().Unix() || err != nil{
			fmt.Println(time.Now().Unix())
			fmt.Println(c.DueTime.Unix())
			return ChannelData{Err:errors.New("error:incorrect date format " )}
		}

		delay,delErr := time.ParseDuration(c.Delay)

		if delErr != nil || delay.Hours() < 0.05{
			return ChannelData{Err:errors.New("error:incorrect delay format")}
		}

		return c
	}else{
		return ChannelData{Err:err}
	}
}


