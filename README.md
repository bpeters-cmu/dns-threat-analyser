# dns-threat-analyser
A golang service that provides graphql endpoints to check if IP addresses are malicious by checking against the spamhause Domain Name System Blocklist
## Prerequisites
For local development the following dependencies are required
```
Go (this project was built using go1.17.2)
sqlite
docker
```
## Getting Started
Clone the project from github
```
git clone https://github.com/bpeters-cmu/dns-threat-analyser.git
```
### Local run instructions
```
cd dns-threat-analyser
go run cmd/server/server.go
```
To change the port, set the `PORT` environment variable before running, by default the server will start on port 8080

### Docker instructions
To change the port the server runs on, modify the `PORT` variable in the below command
```
cd dns-threat-analyser
docker build -t benpeters/dns-threat-analyser .
PORT=8080; docker run -p 127.0.0.1:$PORT:$PORT --env PORT=$PORT benpeters/dns-threat-analyser:latest
```
### How to use the service
After starting the service using the command above it will be available to query at `127.0.0.1:8080/graphql`

The service provides the following queries and mutations:
* mutation enque - Looks up a list of IP's against the spamhause DNSB and stores the result in sqlite. It will return a list with status update for each IP provided
* query getIpDetails - Queries sqlite for a given IP and returns result from DB

The quickest way to get started querying the service is by using the following postman collection [Postman Collection](https://www.getpostman.com/collections/7975261e44b3d5b3673d)

### Examples
Example Query for enque (included in postman collection)
``` 
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
```
variables:
```
{
  "ips": [
    "127.0.0.1", "127.0.0.2", "1.2.3"
  ]
}
```
Sample result:
```
{
    "data": {
        "enque": [
            {
                "error": {
                    "ip_address": "1.2.3",
                    "error_message": "Provided IP is not valid",
                    "error_code": "VALIDATION_ERROR"
                }
            },
            {
                "ip": {
                    "uuid": "8e24298e-56de-49a4-a5ea-356359f01f7a",
                    "created_at": "2021-10-31T19:25:40Z",
                    "updated_at": "2021-10-31T19:39:59Z",
                    "response_code": "NOT LISTED",
                    "ip_address": "127.0.0.1"
                }
            },
            {
                "ip": {
                    "uuid": "c3e91d4b-c9eb-4576-b6db-93d9e196935d",
                    "created_at": "2021-10-31T19:25:40Z",
                    "updated_at": "2021-10-31T19:39:59Z",
                    "response_code": "127.0.0.2\n127.0.0.10\n127.0.0.4",
                    "ip_address": "127.0.0.2"
                }
            }
        ]
    }
}
```
Example query for getIpDetails:
