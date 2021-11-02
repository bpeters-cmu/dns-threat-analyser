package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

// For running on docker
const url string = "http://172.17.0.1:8080/graphql"

// For Local
//const url string = "http://localhost:8080/graphql"

var client *graphql.Client

func main() {
	client = graphql.NewClient(url)
	testAuth()
	testEnque()
	testGetIpDetails()
	log.Println("ALL SMOKE TESTS PASSED")

}
func testAuth() {
	log.Println("---Testing Auth---")
	//Test with invalid credentials
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("wrong", "password")
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if response.Status != "401 Unauthorized" {
		log.Fatal("Basic Auth failed")
	}
	defer response.Body.Close()
	log.Print("---PASS---")

}

func testEnque() {
	log.Print("---Testing enque---")
	enqueRequest := graphql.NewRequest(`
		mutation enque ($ips: [String!]) {
			enque (ips: $ips) {
				... on SuccessStatus {
					ip {
						uuid
						created_at
						updated_at
						response_code
						ip_address
					}
				}
				... on ErrorStatus {
					error {
						ip_address
						error_message
						error_code
					}
				}
			}
		}
	`)

	ctx := context.Background()

	// Map of ip input to expected response_code or error_code
	testCases := map[string]string{
		"127.0.0.1":                              "NOT LISTED",
		"127.0.0.2":                              "127.0.0.2",
		"1.2.3.4":                                "127.0.0.4",
		"1.2.3":                                  "VALIDATION_ERROR",
		"2001:db8:3333:4444:5555:6666:7777:8888": "VALIDATION_ERROR",
		"abc":                                    "VALIDATION_ERROR",
	}

	ips := []string{}
	for key := range testCases {
		ips = append(ips, key)
	}

	enqueRequest.Var("ips", ips)
	enqueRequest.Header.Set("Authorization", "Basic c2VjdXJld29ya3M6c3VwZXJzZWNyZXQ=")

	var enqueResponse map[string][]map[string]map[string]string
	if err := client.Run(ctx, enqueRequest, &enqueResponse); err != nil {
		log.Fatal(err)
	}

	for _, item := range enqueResponse["enque"] {

		if response, ok := item["error"]; ok {
			if !strings.Contains(response["error_code"], testCases[response["ip_address"]]) {
				log.Fatalln("Result does not contain expected response:", response["error_code"], "does not contain:", testCases[response["ip_address"]])
			}

		} else if response, ok := item["ip"]; ok {
			if !strings.Contains(response["response_code"], testCases[response["ip_address"]]) {
				log.Fatalln("Result does not contain expected response:", response["ip_address"], "does not contain:", testCases[response["ip_address"]])
			}
		}
	}
	log.Print("---PASS---")

}

func testGetIpDetails() {
	log.Print("---Testing getIpDetails---")
	getRequest := graphql.NewRequest(`
		query getIPDetails ($ip: String!) {
			getIPDetails (ip: $ip) {
				... on SuccessStatus {
					ip {
						uuid
						created_at
						updated_at
						response_code
						ip_address
					}
				}
				... on ErrorStatus {
					error {
						ip_address
						error_message
						error_code
					}
				}
			}
		}
	`)

	ctx := context.Background()
	getRequest.Header.Set("Authorization", "Basic c2VjdXJld29ya3M6c3VwZXJzZWNyZXQ=")

	// Map of ip input to expected response_code or error_code
	testCases := map[string]string{
		"127.0.0.1":                              "NOT LISTED",
		"127.0.0.2":                              "127.0.0.2",
		"1.2.3.4":                                "127.0.0.4",
		"1.2.3":                                  "VALIDATION_ERROR",
		"2001:db8:3333:4444:5555:6666:7777:8888": "VALIDATION_ERROR",
		"abc":                                    "VALIDATION_ERROR",
	}

	for ip, expectedResult := range testCases {
		getRequest.Var("ip", ip)
		var getIpResponse map[string]map[string]map[string]string
		if err := client.Run(ctx, getRequest, &getIpResponse); err != nil {
			log.Fatal(err)
		}
		if ipDetails, ok := getIpResponse["getIPDetails"]["error"]; ok {
			if !strings.Contains(ipDetails["error_code"], expectedResult) {
				log.Fatalln("Result does not contain expected response:", ipDetails, "does not contain:", expectedResult)
			}

		} else if ipDetails, ok := getIpResponse["getIPDetails"]["ip"]["response_code"]; ok {
			if !strings.Contains(ipDetails, expectedResult) {
				log.Fatalln("Result does not contain expected response:", ipDetails, "does not contain:", expectedResult)
			}
		}

	}
	log.Print("---PASS---")

}
