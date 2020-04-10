package authentication

import (
	"bufio"
	"log"
	"os"
)

type Database struct {
	users map[string]*string
}

func (DB *Database) loadfile() bool {
	var result []string
	file, err := os.Open("./user.db")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return false
	}
	var usr string
	for i, str := range result {
		if i%2 == 0 {
			usr = str
		} else {
			pwd := str
			DB.users[usr] = &pwd
			i = -1
		}
	}
	return true
}
func (DB *Database) Load() {
	DB.users = make(map[string]*string)
	DB.loadfile()
}
func (DB *Database) Verify(username string, password string) bool {
	if *(DB.users[username]) == password {
		return true
	}
	return false
}
