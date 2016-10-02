package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"sort"
	"strings"
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

	fetchMetadata := func(path string, vars ...interface{}) func() (string, error) {
		return func() (string, error) {
			return client.GetMetadata(fmt.Sprintf(path, vars...))
		}
	}

	use := func(value string) func() (string, error) {
		return func() (string, error) {
			return value, nil
		}
	}

	macAddress, err := fetchMetadata("mac")()
	if err != nil {
		logger.Fatal(err)
	}

	vpcIPV4CIDRBlock, err := fetchMetadata("network/interfaces/macs/%s/vpc-ipv4-cidr-block", macAddress)()
	if err != nil {
		logger.Fatal(err)
	}

	calculateDNSServerAddress := func() (string, error) {
		vpcCIDRBlock, _, err := net.ParseCIDR(vpcIPV4CIDRBlock)
		if err != nil {
			return "", err
		}
		vpcDNSServerAddress := vpcCIDRBlock
		incrementIPAddress(vpcDNSServerAddress)
		incrementIPAddress(vpcDNSServerAddress)
		return vpcDNSServerAddress.String(), nil
	}

	fetchSecurityGroups := func() (string, error) {
		value, err := fetchMetadata("security-groups")()
		if err != nil {
			return "", err
		}
		securityGroups := strings.Split(value, "\n")
		sort.Strings(securityGroups)
		return strings.Join(securityGroups, ","), nil
	}

	variables := make(map[string]func() (string, error))
	variables["EC2_ACCOUNT_ID"] = use(metadata.AccountID)
	variables["EC2_ARCHITECTURE"] = use(metadata.Architecture)
	variables["EC2_AVAILABILITY_ZONE"] = use(metadata.AvailabilityZone)
	variables["EC2_IMAGE_ID"] = use(metadata.ImageID)
	variables["EC2_INSTANCE_ID"] = use(metadata.InstanceID)
	variables["EC2_INSTANCE_TYPE"] = use(metadata.InstanceType)
	variables["EC2_KERNEL_ID"] = use(metadata.KernelID)
	variables["EC2_LOCAL_HOSTNAME"] = fetchMetadata("local-hostname")
	variables["EC2_LOCAL_IPV4"] = fetchMetadata("local-ipv4")
	variables["EC2_MAC_ADDRESS"] = use(macAddress)
	variables["EC2_MAC_ADDRESS"] = use(macAddress)
	variables["EC2_PENDING_TIME"] = use(metadata.PendingTime.Format(time.RFC3339))
	variables["EC2_PRIVATE_IP"] = use(metadata.PrivateIP)
	variables["EC2_PUBLIC_HOSTNAME"] = fetchMetadata("public-hostname")
	variables["EC2_PUBLIC_IPV4"] = fetchMetadata("public-ipv4")
	variables["EC2_RAMDISK_ID"] = use(metadata.RamdiskID)
	variables["EC2_REGION"] = use(metadata.Region)
	variables["EC2_RESERVATION_ID"] = fetchMetadata("reservation-id")
	variables["EC2_SECURITY_GROUPS"] = fetchSecurityGroups
	variables["EC2_VPC_DNS_SERVER_ADDRESS"] = calculateDNSServerAddress
	variables["EC2_VPC_ID"] = fetchMetadata("network/interfaces/macs/%s/vpc-id", macAddress)
	variables["EC2_VPC_IPV4_CIDR_BLOCK"] = use(vpcIPV4CIDRBlock)

	stringVariables := make(map[string]string, len(variables))
	for key, value := range variables {
		v, err := value()
		if err != nil {
			logger.Fatal(err)
		}
		stringVariables[key] = v
	}

	keys := make([]string, 0, len(stringVariables))
	for key := range stringVariables {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	output := io.MultiWriter(environmentFile, os.Stderr)
	for _, key := range keys {
		value := stringVariables[key]
		fmt.Fprintf(output, "%s=%s\n", key, value)
	}

}
