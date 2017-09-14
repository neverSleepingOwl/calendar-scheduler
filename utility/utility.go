package utility

//package, containing small functions for different purposes:
//parsing configs, parsing messages, checking time for correctness


type Message struct{
	Title     string	`json:"title"`	//	notification title
	Body      string	`json:"body"`	//	notification body
	IconPath  string	`json:"icon"`	//	path to icon
	DueDate   string	`json:"due"`	//	
	Frequency string	`json:"frequency"`
}
