package domainserver

import (
	""
	"syscall"
	"net"
	"log"
)

//unixSocketServer - struct, responsible for connections via unix domain socket
type unixSocketServer struct{
	Address string `json:"Address"` //	filename of unix socket, get from config
}



//openSocket - function, opening unix socket and waiting for connection
//main unix server loop
func (u unixSocketServer) openSocket(){
	syscall.Unlink(u.Address) //	remove previous socket connections

	l, err := net.Listen("unix", u.Address) //	create socket file and bind it

	if err != nil{
		log.Println("Failed to open unix domain socket ", err.Error())
		return
	}

	log.Println("Unix domain socket bind success, waiting for connection...")

	defer l.Close()	//	close connection on exit

	for{
		fileDescriptor, err := l.Accept()	//	wait for incoming connections
		if err != nil{
			log.Println("Error while accepting connection: ", err.Error())
			return
		}

		log.Println("Domain socket connection accepted, waiting for commands")

		go u.readCommand(fileDescriptor)	//	read commands on accepted connection
	}
}

//readCommand - function, reading commands from a given connection
//puts data to output channel
func (u unixSocketServer)readCommand(connection net.Conn){
	buff:=make([]byte, 1024)

	defer connection.Close()	//	close connection on error

	for{
		if size,err := connection.Read(buff);err == nil{	//	read message

		}else{
			log.Println("Error, incorrect read from unix domain  socket, closing connection")
			return
		}
	}

}

