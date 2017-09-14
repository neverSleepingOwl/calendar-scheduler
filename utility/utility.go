package utility

import (
	"regexp"
	"errors"
	"encoding/json"
	"time"
)

//package, containing small functions for different purposes:
//parsing configs, parsing messages, checking time for correctness


type Notification struct{
	Title     string	`json:"title"`     //	notification title
	Body      string	`json:"body"`      //	notification body
	DueDate   string	`json:"due"`       //	due date
	Delay     string    `json:"delay"`	   //  delay before two notifications
	Frequency uint	    `json:"frequency"` //  how often will notifications appear
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

	command := string(input)

	deleteChecker := regexp.MustCompile(`^(([rR]emove)|([Dd]elete))\s+\w+$`)
	deleteCommand := regexp.MustCompile(`^(([rR]emove)|([Dd]elete))\s+`)
	delimeter:= regexp.MustCompile(`\s`)

	listChecker := "ls"

	if deleteChecker.MatchString(string(input)){
		s := delimeter.ReplaceAllString(command,"")
		return ChannelData{Command:deleteCommand.ReplaceAllString(s,"")}
	}else if command == listChecker{
		return ChannelData{Command:"ls"}
	}

	if err:=json.Unmarshal(input,&c.Notification);err == nil{
		c.DueTime,err = time.Parse(time.RFC822,c.DueDate)

		if c.DueTime.Unix() < time.Now().Unix() || err != nil{
			return ChannelData{Err:errors.New("error:incorrect date format")}
		}

		delay,delErr := time.ParseDuration(c.Delay)

		if delErr != nil || delay.Hours() < 0.1{
			return ChannelData{Err:errors.New("error:incorrect delay format")}
		}

		return c
	}else{
		return ChannelData{Err:err}
	}
}


