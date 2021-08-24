package memory

var Channels map[string](chan bool)

func Setup() {

	Channels = make(map[string]chan bool)

}

func CreateChannel(name string) {

	Channels[name] = make(chan bool)
}

func ActivateChannel(name string) {
	Channels[name] <- true
}
