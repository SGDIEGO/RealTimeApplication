package command

import (
	"strings"
	"sync"
)

var mapCommand = &sync.Map{}

func init() {
	// "/name Diego", name = Diego
	mapCommand.Store("/name", func(name string) string {
		return name
	})
}

func ReadMessage(mssge string) (bool, string) {

	// If string start with '/'
	if mssge[0] == '/' {
		// If mssge = "/name diego", so ind = 5 because index of " " is 5
		ind := strings.Index(mssge, " ")
		option := mssge[:ind]

		command, ok := mapCommand.Load(option)
		if !ok || (command == nil) {
			return false, ""
		}

		// Parse to func(name string)
		commandFunc, ok := command.(func(name string) string)
		if !ok {
			return false, ""
		}
		// Pass rest of string
		return true, commandFunc(mssge[ind+1:])
	}

	return false, mssge
}
