package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	environmentFilePath = "/etc/ec2/environment"
)

func incrementIPAddress(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	logger := log.New(os.Stderr, "", 0)

	sess, err := session.NewSession()
	if err != nil {
		logger.Fatal(err)
	}

	client := ec2metadata.New(sess)
	if !client.Available() {
		logger.Fatal("ec2 metadata service unavailable")
	}

	err = os.Mkdir(path.Dir(environmentFilePath), 0775)
	if err != nil {
		if !os.IsExist(err) {
			logger.Fatal(err)
		}
	}

	environmentFile, err := os.Create(environmentFilePath)
	if err != nil {
		logger.Fatal(err)
	}

	metadata, err := client.GetInstanceIdentityDocument()
	if err != nil {
		logger.Fatal(err)
	}

	macAddress, err := client.GetMetadata("mac")
	if err != nil {
		logger.Fatal(err)
	}

	vpcID, err := client.GetMetadata(fmt.Sprintf("network/interfaces/macs/%s/vpc-id", macAddress))
	if err != nil {
		logger.Fatal(err)
	}

	vpcIPV4CIDRBlock, err := client.GetMetadata(fmt.Sprintf("network/interfaces/macs/%s/vpc-ipv4-cidr-block", macAddress))
	if err != nil {
		logger.Fatal(err)
	}

	vpcCIDRBlock, _, err := net.ParseCIDR(vpcIPV4CIDRBlock)
	if err != nil {
		logger.Fatal(err)
	}
	vpcDNSServerAddress := vpcCIDRBlock
	incrementIPAddress(vpcDNSServerAddress)
	incrementIPAddress(vpcDNSServerAddress)

	output := io.MultiWriter(environmentFile, os.Stderr)
	writeEnvironmentVariable := func(name, value string) {
		fmt.Fprintf(output, "%s=%s\n", name, value)
	}

	writeEnvironmentVariable("EC2_ACCOUNT_ID", metadata.AccountID)
	writeEnvironmentVariable("EC2_ARCHITECTURE", metadata.Architecture)
	writeEnvironmentVariable("EC2_AVAILABILITY_ZONE", metadata.AvailabilityZone)
	writeEnvironmentVariable("EC2_IMAGE_ID", metadata.ImageID)
	writeEnvironmentVariable("EC2_INSTANCE_ID", metadata.InstanceID)
	writeEnvironmentVariable("EC2_INSTANCE_TYPE", metadata.InstanceType)
	writeEnvironmentVariable("EC2_KERNEL_ID", metadata.KernelID)
	writeEnvironmentVariable("EC2_MAC_ADDRESS", macAddress)
	writeEnvironmentVariable("EC2_PENDING_TIME", metadata.PendingTime.Format(time.RFC3339))
	writeEnvironmentVariable("EC2_PRIVATE_IP", metadata.PrivateIP)
	writeEnvironmentVariable("EC2_RAMDISK_ID", metadata.RamdiskID)
	writeEnvironmentVariable("EC2_REGION", metadata.Region)
	writeEnvironmentVariable("EC2_VPC_ID", vpcID)
	writeEnvironmentVariable("EC2_VPC_IPV4_CIDR_BLOCK", vpcIPV4CIDRBlock)
	writeEnvironmentVariable("EC2_VPC_DNS_SERVER_ADDRESS", vpcDNSServerAddress.String())
}
