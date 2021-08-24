package util

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Error : Respuesta segun formato de error estandar
type Error struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func GetError(message string, err error) Error {
	if err != nil {
		return Error{Message: message, Error: err.Error()}
	}
	return Error{Message: message}

}
