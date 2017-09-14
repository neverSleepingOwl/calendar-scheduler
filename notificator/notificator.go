package notificator

import "scheduler/utility"

type Notifier struct{
	Done chan struct{}
	utility.Notification
	
}