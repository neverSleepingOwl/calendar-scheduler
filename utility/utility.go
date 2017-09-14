package utility

import (
	"regexp"
	"strings"
	"errors"
)

//package, containing small functions for different purposes:
//parsing configs, parsing messages, checking time for correctness


type Notification struct{
	Title     string	//	notification title
	Body      string	//	notification body
	IconPath  string	//	path to icon
	DueDate   uint32	//	due date
	Frequency uint32	//  how often will notifications appear
}

//struct for channel output
type ChannelData struct{
	Notification
	Err error
}

//NweNotification parse message
func NewNotification(input []byte)(c ChannelData){
	dataToParse := string(input)

	checker := regexp.MustCompile(`^(\w+:[^;\s]+\s*;\s*){2,5}(\w+:[^;\s]+)$`)
	titleChecker := regexp.MustCompile(`^\s*[tT]itle\s*:\s*\w+$`)
	bodyChecker := regexp.MustCompile(`^\s*([bB]ody)|([Mm]ess(a|[ea]n)ge)\s*:\s*.+$`)
	iconPathChecker := titleChecker
	dueDateChecker := regexp.MustCompile(`^\s*[Dd]ue\s*:.+`)
	freqChecker := regexp.MustCompile(`^\s*[fF]requency\s*:\s*[1-9]+\d\s*(h|min|m|d(ay[s]?)?)\s*$`)
	if !checker.MatchString(dataToParse){
		return ChannelData{Notification{}, errors.New("incorrect input, syntax error")}
	}
	splited := strings.Split(dataToParse," ")
	for _,element:= range splited{
		switch{
		case titleChecker.MatchString(element):
			if title:=strings.Split(trimWhitespaces(element),":"); len(title) > 1{
				c.Title = title[1]
			}else{

			}
		case bodyChecker.MatchString(element):
			if body:=strings.Split(element,":"); len(body) > 1{
				c.Body = body[1]
			}
		case iconPathChecker.MatchString(element):
			if icon:=strings.Split(trimWhitespaces(element),":"); len(icon) > 1{
				c.IconPath= icon[1]
			}
		case dueDateChecker.MatchString(element):
		case freqChecker.MatchString(element):
		}
	}
}

func trimWhitespaces(s string)string{
	return strings.Replace(s, " ", "",-1)
}

func CheckDate(date string)bool{
	dateChecker := regexp.MustCompile(`^([0-9]?\d\.){2}[1-9]\d{3}()$`)


}



