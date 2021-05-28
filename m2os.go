/*
* jmu april 2021
*/
package main

import (
	"fmt"
	"os"
	"log"
	"context"
	"github.com/google/uuid"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/streaming"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// reads from mosquito and writes in oci streaming
func main() {

    intro()

	broker := os.Getenv("broker")
    fmt.Printf("mqtt endpoint: %s\n", broker)

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

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(uuid.New().String())
	qos := 0
	log.Print(opts)

	nrcp := common.NewRawConfigurationProvider(tenancy, user, region, fingerprint, string(ppkcontent), &ppkpasswd)

	var endpoint = "https://cell-1.streaming." + region + ".oci.oraclecloud.com"
	fmt.Printf("Streaming endpoint: %s\n", endpoint)
	sclient, err := streaming.NewStreamClientWithConfigurationProvider(nrcp, endpoint)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
		receiveCount := 0
		choke := make(chan [2]string)

		opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		})

		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := client.Subscribe(mqtttopic, byte(qos), nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		for receiveCount >=0  {
			incoming := <-choke
			fmt.Printf("Recieved from mqtt topic %s the message <- %s\t", incoming[0], incoming[1])
			receiveCount++
			putMessage(sclient, stream, region, uuid.New().String(), fmt.Sprintf(incoming[1]))
		}

		client.Disconnect(250)
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
	fmt.Println("Sent to OCI streaming -> \n", value)
	return  0
}
/*
*
*/
func intro(){
    fmt.Fprintf(os.Stderr, "\n (c) jmu 2021 | m2os, get messages from mosquitto topic and publishes to oci streaming\n")
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