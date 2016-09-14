package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const ami = "ami-2d39803a"

var wg sync.WaitGroup

func launchInstances(count int, svc *ec2.EC2) *string {

	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(ami),
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

func fetchRunningReservation(reservationID *string, svc *ec2.EC2) *ec2.Reservation {

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
	print("Waiting for instance(s) to launch")
	for {

		// keep the user notified of progress
		print(".")
		time.Sleep(time.Second)

		// fetch the instances
		reservations, err := svc.DescribeInstances(params)
		if err != nil {
			panic(err)
		}

		// check if all the instances in the reservation are running
		allRunning := true
		reservation := reservations.Reservations[0]
		for _, instance := range reservation.Instances {
			allRunning = allRunning && (*instance.State.Name == "running")
		}
		if allRunning {
			return reservation
		}
	}
}

func prepareSSHclient() *ssh.ClientConfig {
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
	return config
}

func installCassandra(instance *ec2.Instance) {
	defer wg.Done()
	for {
		conn, err := ssh.Dial("tcp", *instance.PublicDnsName+":22", prepareSSHclient())

		if err != nil {
			switch t := err.(type) {
			case *net.OpError:
				log.Println(t.Error())
				time.Sleep(5 * time.Second)
				continue
			default:
				panic(err)
			}
		}
		log.Println("Connected to", *instance.PublicDnsName)
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()

		var stdout bytes.Buffer
		session.Stdout = &stdout
		addDependency := "printf '\n" +
			"deb http://www.apache.org/dist/cassandra/debian 21x main\n" +
			"deb-src http://www.apache.org/dist/cassandra/debian 21x main' | " +
			"sudo tee -a /etc/apt/sources.list"
		addKey := "gpg --keyserver pgp.mit.edu --recv-keys 749D6EEC0353B12C && " +
			"gpg --export --armor 749D6EEC0353B12C | " +
			"sudo apt-key add -"
		installCassandra := "sudo apt-get update && sudo apt-get -y install cassandra"
		session.Run(addDependency + " && " + addKey + " && " + installCassandra)

		log.Println(stdout.String())
		return
	}
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

	// fetch running reservation
	reservation := fetchRunningReservation(reservationID, svc)

	// install cassandra on all the instances
	for _, instance := range reservation.Instances {
		wg.Add(1)
		go installCassandra(instance)
	}
	wg.Wait()
}
