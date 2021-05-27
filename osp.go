/*
* jmu may 2021, oracle streaming random producer
*/
package main

import (
	"fmt"
	"os"
	"context"
	"github.com/google/uuid"
	"math/rand"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/streaming"
)

// reads from mosquito and writes in oci streaming
func main() {

    intro()

    mqtttopic := os.Getenv("mqtttopic")
    fmt.Printf("mqtt topic name: %s\n", mqtttopic)

	stream := os.Getenv("stream")
    fmt.Printf("oci stream ocid: %s\n", stream)

	tenancy := os.Getenv("tenancy")
    fmt.Printf("oci tenancy ocid: %s\n", tenancy)

    user := os.Getenv("user")
    fmt.Printf("oci user ocid: %s\n", user)

    region := os.Getenv("region")
    fmt.Printf("oci region: %s\n", region)

    fingerprint := os.Getenv("fingerprint")
    fmt.Printf("fingerprint: %s\n", fingerprint)

    ppkcontent := os.Getenv("ppkfile")

	ppkpasswd := os.Getenv("ppkpassword")
    fmt.Printf("ppk passwd: %s\n", "****")

	nrcp := common.NewRawConfigurationProvider(tenancy, user, region, fingerprint, string(ppkcontent), &ppkpasswd)

	var endpoint = "https://cell-1.streaming." + region + ".oci.oraclecloud.com"
	fmt.Printf("Streaming endpoint: %s\n", endpoint)
	sclient, err := streaming.NewStreamClientWithConfigurationProvider(nrcp, endpoint)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var message = ""

	for 1 >=0  {
		message = fmt.Sprintf("{t:%d,d:%d,v:%d}", rand.Intn(9876543210000000), rand.Intn(1234567890000000), rand.Intn(90817263000000)) 
		putMessage(sclient, stream, region, uuid.New().String(), message)
	}

	fmt.Println("Sample Subscriber Disconnected")
}

/*
* put a meessage in a topic
*/
func putMessage(client streaming.StreamClient, stream string, region string, key string, value string) int{
	var req streaming.PutMessagesRequest
	req.StreamId = &stream
	var entry streaming.PutMessagesDetailsEntry
	entry.Key = []byte(key)
	entry.Value = []byte(value)
	var entryarray [] streaming.PutMessagesDetailsEntry
	entryarray = append(entryarray, entry)
	var det streaming.PutMessagesDetails
	det.Messages = entryarray
	req.PutMessagesDetails = det
	//client.SetRegion(region)
	_, err := client.PutMessages(context.Background(), req)
	if err != nil {
		fmt.Println("Error: ", err)
		return  -1
	}
	fmt.Println("Sent to OCI streaming -> ", value)
	return  0
}
/*
*
*/
func intro(){
    fmt.Fprintf(os.Stderr, "\n (c) jmu 2021 | osp, produces random messages and publishes to oci streaming\n")
    fmt.Printf("--------------------------------------------------------------------------------------\n")
}
/*
*
*/
func check(e error) {
    if e != nil {
        panic(e)
    }
}