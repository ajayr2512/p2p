package peers

import (
	"bufio"
	"bytes"
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
		// To DO: If empty then also -1 ??
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
			cookie = -1
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

func (peer *Peer) Register() (err error) {
	// TO DO: Use protoBUf ????
	conn, err := net.Dial("tcp", "localhost:60000")
	if err != nil {
		return err
	}
	var msg []byte
	if msg, err = RegisterRequest(*peer); err != nil {
		return err
	}
	log.Println(string(msg))
	if _, err = conn.Write(msg); err != nil {
		return err
	}

	b := bufio.NewReader(conn)
	var req []byte
	if req, err = b.ReadBytes('\r'); err != nil {
		return err
	}

	br := bytes.NewBuffer(req)
	scanner := bufio.NewScanner(br)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		if s[0] == "STATUS" && s[1] != "NEW" {
			break
		}
		if s[0] == "COOKIE" {
			cookiePath := filepath.Join(peer.data, "COOKIE")
			file, err := os.OpenFile(cookiePath, os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err = file.WriteString(s[1]); err != nil {
				return err
			}
			peer.cookie, _ = strconv.Atoi(s[1])
		}
	}
	return nil

}

func (peer *Peer) Leave() (err error) {
	conn, err := net.Dial("tcp", "localhost:60000")
	if err != nil {
		return err
	}
	var msg []byte
	if msg, err = LeaveRequest(*peer); err != nil {
		return err
	}
	log.Println(string(msg))
	if _, err = conn.Write(msg); err != nil {
		return err
	}

	if _, err = conn.Write(msg); err != nil {
		return err
	}
	return nil
}
