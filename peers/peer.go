package peers

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Peer struct {
	hostName string
	cookie   int
	port     int
	data     string
}

type fileList []string

func NewPeer(port int, hname, data string) (peer Peer, files fileList) {
	// create cookie file if doesnt exist
	cookie := -1
	cookiePath := filepath.Join(data, "COOKIE")
	var file *os.File
	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		log.Println("Here 1")
		if file, err = os.Create(cookiePath); err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		log.Println("Here 2")
		file, err := os.OpenFile(cookiePath, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		var cookieBytes []byte
		if cookieBytes, err = ioutil.ReadAll(file); err != nil {
			log.Fatal(err)
		}
		if cookie, err = strconv.Atoi(string(cookieBytes)); err != nil {
			log.Println("Corrupt cookie file")
			log.Println(err)
		}
		defer file.Close()

	}

	log.Println("cookie: ", cookie)
	// data directory check to create file list
	fileinfo, err := ioutil.ReadDir(data)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fileinfo {
		files = append(files, file.Name())
	}
	return Peer{
		hostName: hname,
		port:     port,
		data:     data,
		cookie:   cookie,
	}, files

}

func (peer Peer) Register() (err error) {
	// TO DO: Use protoBUf ????
	conn, err := net.Dial("tcp", "localhost:60000")
	if err != nil {
		return err
	}
	var msg []byte
	if msg, err = RegisterMessage(peer); err != nil {
		return err
	}
	if _, err = conn.Write(msg); err != nil {
		return err
	}
	msg, err = ioutil.ReadAll(conn)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(msg)
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		switch s[0] {
		case "STATUS":
			if s[1] == "FAIL" {
				return errors.New("protocl error")
			}
		case "COOKIE":
			// err handle
			peer.cookie, _ = strconv.Atoi(s[1])
		default:
			log.Printf("Received Unknown info from RS: %s", s[0])
		}
	}
	return nil

}
