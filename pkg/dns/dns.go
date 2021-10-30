package dns

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const spamhause = "zen.spamhaus.org"

func HandleDnsLookups(ips []string) error {
	log.Println("handle dns lookups")
	for _, ip := range ips {
		reversedIp := reverse(strings.Split(ip, "."))
		query := strings.Join(reversedIp, ".") + "." + spamhause
		log.Println(query)
		stdout, err := exec.Command("dig", "+short", query).Output()
		if err != nil {
			log.Println("Err:" + err.Error())
			return errors.New(fmt.Sprint("ERROR executing host command for dns lookup:", ip))
		}
		log.Println(string(stdout))

	}
	return nil
}

//referenced from https://golangcookbook.com/chapters/arrays/reverse/
func reverse(s []string) []string {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}
