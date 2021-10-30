package dns

import (
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/bpeters-cmu/dns-threat-analyser/graph/model"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/database"
	"github.com/google/uuid"
)

const spamhauseUrl = ".zen.spamhaus.org"

func HandleDnsLookup(ipAddr string, resultsChan chan model.EnqueStatus) {
	// Validate IP first
	if net.ParseIP(ipAddr) == nil {
		resultsChan <- model.EnqueError{IP: ipAddr, Message: "Provided IP is not valid"}
		return
	}
	if strings.Contains(ipAddr, ":") {
		resultsChan <- model.EnqueError{IP: ipAddr, Message: "IPv6 is not supported"}
		return
	}
	// Format dns lookup query
	reversedIp := reverse(strings.Split(ipAddr, "."))
	query := strings.Join(reversedIp, ".") + spamhauseUrl
	// Call spamhause url for lookup
	stdout, err := exec.Command("dig", "+short", query).Output()
	if err != nil {
		resultsChan <- model.EnqueError{IP: ipAddr, Message: "ERROR executing host command for dns lookup"}
		return
	}
	// Check if IP already exist in DB
	db := database.SqliteDB{}
	existingIp, err := db.GetIp(ipAddr)

	if err != nil {
		// IP doesn't exist in DB so will attempt to create and insert a new IP
		ip := model.IP{IPAddress: ipAddr, UUID: uuid.NewString(), ResponseCode: string(stdout), CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if err = db.SaveIp(&ip); err != nil {
			resultsChan <- model.EnqueError{IP: ipAddr, Message: "ERROR saving IP to database"}
			return
		}
		resultsChan <- model.EnqueSuccess{IP: &ip}
		return
	}
	//IP is already present in DB so we will update it
	existingIp.ResponseCode = string(stdout)
	existingIp.UpdatedAt = time.Now()
	if err = db.SaveIp(existingIp); err != nil {
		resultsChan <- model.EnqueError{IP: ipAddr, Message: "ERROR saving IP to database"}
		return
	}

	resultsChan <- model.EnqueSuccess{IP: existingIp}
	return

}

//referenced from https://golangcookbook.com/chapters/arrays/reverse/
func reverse(s []string) []string {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}
