package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const AMI = "ami-2d39803a"

func launchInstances(count int, svc *ec2.EC2) *string {

	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(AMI),
		MaxCount:     aws.Int64(int64(count)),
		MinCount:     aws.Int64(int64(count)),
		InstanceType: aws.String("t2.micro"),
		KeyName:      aws.String("testing"),
	}
	reservation, err := svc.RunInstances(params)
	if err != nil {
		panic(err)
	}

	return reservation.ReservationId
}

func main() {

	// take in size of cluster as command line argument
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	// start session
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		panic(err)
	}
	svc := ec2.New(sess)

	// launch instance
	reservationID := launchInstances(count, svc)

	// prepare to fetch instance information
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("reservation-id"),
				Values: []*string{
					reservationID,
				},
			},
		},
	}

	// start attempting to fetch instances
	reservations := new(ec2.DescribeInstancesOutput)
	print("Waiting for instance(s) to launch")
	for {

		// keep the user notified of progress
		print(".")
		time.Sleep(time.Second)

		// fetch the instances
		reservations, err = svc.DescribeInstances(params)
		if err != nil {
			panic(err)
		}

		// check if all the instances in the reservation are running
		allRunning := true
		for _, reservation := range reservations.Reservations {
			for _, instance := range reservation.Instances {
				allRunning = allRunning && (*instance.State.Name == "running")
			}
		}
		if allRunning {
			break
		}
	}
	println()

	// prepare SSH client
	key, err := ioutil.ReadFile("/Users/plato/Desktop/testing.pem")
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	auths := []ssh.AuthMethod{ssh.PublicKeys(signer)}
	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: auths,
	}

	// get the the public domains
	for _, reservation := range reservations.Reservations {
		for _, instance := range reservation.Instances {
			conn, err := ssh.Dial("tcp", *instance.PublicDnsName+":22", config)

			// need to handle a refused connection here
			if err != nil {
				panic(err)
			}
			session, err := conn.NewSession()
			if err != nil {
				panic(err)
			}
			defer session.Close()

			var stdout bytes.Buffer
			session.Stdout = &stdout
			session.Run("printf '\ndeb http://www.apache.org/dist/cassandra/debian 21x main\ndeb-src http://www.apache.org/dist/cassandra/debian 21x main' | sudo tee -a /etc/apt/sources.list && gpg --keyserver pgp.mit.edu --recv-keys 749D6EEC0353B12C && gpg --export --armor 749D6EEC0353B12C | sudo apt-key add - && sudo apt-get update && sudo apt-get -y install cassandra")

			fmt.Println(stdout.String())
		}
	}

}
