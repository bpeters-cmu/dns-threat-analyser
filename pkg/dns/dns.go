package dns

import (
	"errors"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/bpeters-cmu/dns-threat-analyser/graph/model"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/database"
	"github.com/google/uuid"
)

const (
	spamhauseUrl    = ".zen.spamhaus.org"
	ValidationError = "VALIDATION_ERROR"
	SystemError     = "SYSTEM_ERROR"
)

type DnsClient interface {
	Dig(query string) (string, error)
}

type DigClient struct {
}

// Dig sends a dig command to the spamhause service and formats the response code
func (dc *DigClient) Dig(query string) (string, error) {
	stdout, err := exec.Command("dig", "+short", query).Output()
	if err != nil {
		return "", errors.New("ERROR executing host command for dns lookup")
	}
	// an ip lookup can have multiple return codes each on a new line
	// this will clean it up to be more readable
	responseCode := string(stdout)
	responseCode = strings.TrimRight(responseCode, "\n")
	responseCode = strings.ReplaceAll(responseCode, "\n", ", ")
	// Spamhause returns empty string if record is not listed for dig +short cmd
	if responseCode == "" {
		responseCode = "NOT LISTED"
	}

	return responseCode, nil
}

// HandleDnsLookup validates the input ip, calls the spamhause service
// for the lookup then creates or updates a db entry for the IP
func HandleDnsLookup(ipAddr string, db database.Database, dc DnsClient, resultsChan chan model.Status) {
	// Validate IP before lookup
	if err := ValidateIp(ipAddr); err != nil {
		resultsChan <- model.ErrorStatus{Error: &model.Error{IPAddress: ipAddr, ErrorMessage: err.Error(), ErrorCode: ValidationError}}
		return
	}
	// Format dns lookup query
	reversedIp := reverse(strings.Split(ipAddr, "."))
	query := strings.Join(reversedIp, ".") + spamhauseUrl
	// Call spamhause url for lookup
	responseCode, err := dc.Dig(query)
	if err != nil {
		resultsChan <- model.ErrorStatus{Error: &model.Error{IPAddress: ipAddr, ErrorMessage: err.Error(), ErrorCode: SystemError}}
		return
	}

	// Check if IP already exist in DB
	existingIp, err := db.GetIp(ipAddr)
	if err != nil {
		resultsChan <- model.ErrorStatus{Error: &model.Error{IPAddress: ipAddr, ErrorMessage: err.Error(), ErrorCode: SystemError}}
		return
	}
	if existingIp == nil {
		// IP doesn't exist in DB so will attempt to create and insert a new IP
		ip := model.IP{IPAddress: ipAddr, UUID: uuid.NewString(), ResponseCode: responseCode, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if err = db.SaveIp(&ip); err != nil {
			resultsChan <- model.ErrorStatus{Error: &model.Error{IPAddress: ipAddr, ErrorMessage: err.Error(), ErrorCode: SystemError}}
			return
		}
		resultsChan <- model.SuccessStatus{IP: &ip}
		return
	}
	//IP is already present in DB so we will update it
	existingIp.ResponseCode = responseCode
	existingIp.UpdatedAt = time.Now()
	if err = db.SaveIp(existingIp); err != nil {
		resultsChan <- model.ErrorStatus{Error: &model.Error{IPAddress: ipAddr, ErrorMessage: err.Error(), ErrorCode: SystemError}}
		return
	}

	resultsChan <- model.SuccessStatus{IP: existingIp}
	return

}

// ValidateIp validates that the input string
// is a valid IPv4 address
func ValidateIp(ipAddr string) error {
	if net.ParseIP(ipAddr) == nil {
		return errors.New("Provided IP is not valid")
	}
	if strings.Contains(ipAddr, ":") {
		return errors.New("IPv6 is not supported")
	}
	return nil
}

//reverse a string array, referenced from https://golangcookbook.com/chapters/arrays/reverse/
func reverse(s []string) []string {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}
