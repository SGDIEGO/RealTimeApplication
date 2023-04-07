package server

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
)

type client struct {
	conn       net.Conn
	idClient   int
	nameClient string
}

func Client(conn net.Conn) *client {
	return &client{
		conn:     conn,
		idClient: rand.Int(),
	}
}

func (c *client) choose(opt string, value string) (bool, string) {
	// "/name Diego", name = Diego
	value = value[:len(value)-2]
	switch opt {
	case "/name":
		return true, func(name string) string {
			var response string

			if c.nameClient == "" {
				response = fmt.Sprintf("User (%v) change name to %v", c.idClient, name)
			} else {
				response = fmt.Sprintf("User (%v) change name to %v", c.nameClient, name)
			}

			c.nameClient = name
			return response
		}(value)

	default:
		return false, ""
	}
}

func (c *client) ReadMssage(mssge string) (bool, string) {
	// If string start with '/'
	if mssge[0] == '/' {
		// If mssge = "/name diego", so ind = 5 because index of " " is 5
		ind := strings.Index(mssge, " ")
		option := mssge[:ind]

		// Pass rest of string
		commandValid, mssgeCommand := c.choose(option, mssge[ind+1:])
		if !commandValid {
			return true, ""
		}

		return true, mssgeCommand
	}

	return false, mssge
}

func (c *client) showMessage(mssge string) string {
	if c.nameClient == "" {
		return fmt.Sprintf("%v:%v", c.idClient, mssge)
	}
	return fmt.Sprintf("%v:%v", c.nameClient, mssge)
}

func (c *client) readerMssge(connMap *sync.Map) {
	defer func() {
		c.conn.Close()             // Close connection for user
		connMap.Delete(c.idClient) // Remove idClient from map
	}()

	for {
		// Read message
		mssge, err := bufio.NewReader(c.conn).ReadString('\n')
		isCommand, newMssge := c.ReadMssage(mssge)
		if isCommand {
			if newMssge == "" {
				c.conn.Write([]byte("command invalid!"))
				continue
			}
			mssge = newMssge
		}
		if err != nil {
			log.Println(err)
			return
		}

		// Show message to all users before connected
		connMap.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if _, err := conn.Write([]byte(c.showMessage(mssge))); err != nil {
					log.Fatal("error on writing to connection", err)
				}
			}
			return true
		})

	}
}

func (c *client) ListenMssge(connMap *sync.Map) {
	// Show connection information
	log.Println("New connection: ", c.conn.RemoteAddr().String())

	// Show new connection to all users
	connMap.Range(func(key, value interface{}) bool {
		if conn, ok := value.(net.Conn); ok {
			if _, err := conn.Write([]byte(fmt.Sprintln("New user connected: ", c.idClient))); err != nil {
				log.Fatal("error on writing to connection", err)
			}
		}
		return true
	})

	// Save new connection
	connMap.Store(c.idClient, c.conn)

	// Listenning messages from the new connection
	go c.readerMssge(connMap)
}
