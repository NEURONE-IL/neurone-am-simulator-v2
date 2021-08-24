package hash

var Channels map[string](chan bool)

func Septup() {

	Channels = make(map[string]chan bool)

}
