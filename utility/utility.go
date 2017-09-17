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

//{"title":"tag","body":"testing arp requests","due":"17 Sep 17 14:30 GMT","delay":"0.05h","icon":"/usr/share/icons/gnome/32x32/emotes/face-monkey.png"}

type Notification struct{
	Title     string	`json:"title"`     //	notification title
	Body      string	`json:"body"`      //	notification body
	DueDate   string	`json:"due"`       //	due date
	Delay     string    `json:"delay"`	   //  delay before two notifications
	Icon 	  string 	`json:"icon"`
	DueTime   time.Time                    //	due date, but has time format
}

//struct for channel output
type ChannelData struct{
	Notification
	Command string
	Err error
}

//NewNotification parse message
//Message can be either a notification in json format or
//command like 'ls', 'js' 'delete|remove <notification>'
//function firstly filters commands and returns notification with set command field
//if message isn't command, check fields for validity and then return notification with given parameters

//commands are:delete, ls and js
//delete: syntax - delete or remove <record title>, removes given notification
//ls: syntax - ls, lists all notifications
//js: lists all notifications in json format, command required for easy gui implementation
func NewNotification(input []byte)(ChannelData) {
	var c ChannelData = ChannelData{Notification:Notification{},Err:nil}
	//if string is \n terminated remove the last character, todo fix
	if string(input)[len(string(input))-1] == '\n'{
		input = input[:len(input)-1]
	}
	command := string(input)

	//regexp to recognise delete command
	deleteChecker := regexp.MustCompile(`^(([rR]emove)|([Dd]elete))\s+\w+$`)
	//regexp to remove 'delete  ' part of delete command and send just key in hashtable as a command
	//maybe too symplistic, because we can't use titles like ls or js, so todo probably fix
	deleteCommand := regexp.MustCompile(`^(([rR]emove)|([Dd]elete))`)
	//we need to remove all spaces from delete command to leave just title of record we should delete
	delimeter:= regexp.MustCompile(`\s`)

	listChecker := "ls"	//	ls command checker
	jsonChecker := "js" // 	json command checker


	// if command is delete command
	if deleteChecker.MatchString(string(input)){

		s := delimeter.ReplaceAllString(command,"")		// trim all whitespaces
		s = deleteCommand.ReplaceAllString(s,"")		// remove delete part of command

		return ChannelData{Command:s}	//	send name of record to delete as command

	}else if command == listChecker || command == jsonChecker{		//	if it's simple command like ls or js
		return ChannelData{Command:command}		//
	}

	//if notification is sent
	if err:=json.Unmarshal(input,&c.Notification);err == nil{
		c.DueTime,err = time.Parse(time.RFC822,c.DueDate)	//	parse time

		if c.DueTime.Unix() < time.Now().Unix() || err != nil{	//	check time for validity
			return ChannelData{Err:errors.New("error:incorrect date format\n" )}
		}

		delay,delErr := time.ParseDuration(c.Delay)

		if delErr != nil || delay.Hours() < 0.05{	//	check delays for validity
			return ChannelData{Err:errors.New("error:incorrect delay format\n")}
		}

		return c
	}else{
		return ChannelData{Err:err}
	}
}


