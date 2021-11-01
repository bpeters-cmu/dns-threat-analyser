package dns

import (
	"errors"
	"testing"
	"time"

	"github.com/bpeters-cmu/dns-threat-analyser/graph/model"
	"github.com/google/uuid"
)

// Mock for DB
var saveIpMock func(ip *model.IP) error
var getIpMock func(pAddr string) (*model.IP, error)

type MockDB struct{}

func (db *MockDB) SaveIp(ip *model.IP) error {
	return saveIpMock(ip)
}

func (db *MockDB) GetIp(ipAddr string) (*model.IP, error) {
	return getIpMock(ipAddr)
}

// Mock for dig client
var digMock func(query string) (string, error)

type MockDigClient struct{}

func (dc *MockDigClient) Dig(query string) (string, error) {
	return digMock(query)
}
func TestHandleDnsLookup(t *testing.T) {

	saveIpMock = func(ip *model.IP) error {
		return nil
	}

	// returning nil, nil hear means ip does not exist in database
	getIpMock = func(ipAddr string) (*model.IP, error) {
		return nil, nil
	}

	digMock = func(query string) (string, error) {
		return "127.0.0.2", nil
	}

	resultsChan := make(chan model.Status, 1)

	HandleDnsLookup("1.2.3.4", &MockDB{}, &MockDigClient{}, resultsChan)

	result := <-resultsChan

	switch result.(type) {
	case model.SuccessStatus:
		return
	default:
		t.Fail()
	}

}

func TestHandleDnsLookupForExisting(t *testing.T) {
	saveIpMock = func(ip *model.IP) error {
		return nil
	}

	currentTime := time.Now()
	// returning nil, nil hear means ip does not exist in database
	getIpMock = func(ipAddr string) (*model.IP, error) {
		return &model.IP{IPAddress: "1.2.3.4", UUID: uuid.NewString(), CreatedAt: currentTime, UpdatedAt: currentTime}, nil
	}

	digMock = func(query string) (string, error) {
		return "127.0.0.2", nil
	}

	resultsChan := make(chan model.Status, 1)

	HandleDnsLookup("1.2.3.4", &MockDB{}, &MockDigClient{}, resultsChan)

	result := <-resultsChan

	switch result.(type) {
	case model.SuccessStatus:
		success := result.(model.SuccessStatus)
		// Updated time should be after created time
		if !success.IP.UpdatedAt.After(currentTime) {
			t.Fail()
		}
	default:
		t.Fail()
	}

}

func TestHandleDnsLookupForInvalidIp(t *testing.T) {
	saveIpMock = func(ip *model.IP) error {
		return nil
	}

	// returning nil, nil hear means ip does not exist in database
	getIpMock = func(ipAddr string) (*model.IP, error) {
		return nil, nil
	}

	digMock = func(query string) (string, error) {
		return "", nil
	}

	resultsChan := make(chan model.Status, 1)

	HandleDnsLookup("1.2.3", &MockDB{}, &MockDigClient{}, resultsChan)

	result := <-resultsChan

	switch result.(type) {
	case model.ErrorStatus:
		// Updated time should be after created time
		errorCode := result.(model.ErrorStatus).Error.ErrorCode
		if errorCode != ValidationError {
			t.Fail()
		}
	default:
		t.Fail()
	}

}

func TestHandleDnsLookupWithDBWriteFailure(t *testing.T) {
	saveIpMock = func(ip *model.IP) error {
		return errors.New("WRITE ERROR")
	}

	// returning nil, nil hear means ip does not exist in database
	getIpMock = func(ipAddr string) (*model.IP, error) {
		return nil, nil
	}

	digMock = func(query string) (string, error) {
		return "", nil
	}

	resultsChan := make(chan model.Status, 1)

	HandleDnsLookup("1.2.3.4", &MockDB{}, &MockDigClient{}, resultsChan)

	result := <-resultsChan

	switch result.(type) {
	case model.ErrorStatus:
		// Updated time should be after created time
		errorCode := result.(model.ErrorStatus).Error.ErrorCode
		if errorCode != SystemError {
			t.Fail()
		}
	default:
		t.Fail()
	}

}

func TestValidateIp(t *testing.T) {
	validIps := []string{"127.0.0.1", "1.2.3.4", "0.0.0.0"}
	for _, ip := range validIps {
		err := ValidateIp(ip)
		if err != nil {
			t.Fail()
		}
	}

	invalidIps := []string{"abc", "2001:db8:3333:4444:5555:6666:7777:8888", "1.2.3", ""}

	for _, ip := range invalidIps {
		err := ValidateIp(ip)
		if err == nil {
			t.Fail()
		}
	}

}
