package cmd

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Starts a new AWS session
func startAWSSession(IP string, Machine string) {
	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := ec2.New(session, &aws.Config{Region: aws.String(Region)})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(Machine),
				},
			},
		},
	}
	response, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	for _, inst := range response.Reservations[0].Instances {
		for _, securityGroup := range inst.SecurityGroups {

			ipExists := IPExistsOnSecurityGroup(IP, svc, securityGroup)
			if !ipExists {
				AddIPAddressOnSecurityGroup(svc, securityGroup.GroupId)
			} else {
				log.Print("WARNING: This IP already has access!")
			}
		}
	}

}

// Verify is IP exists on SecurityGroup
func IPExistsOnSecurityGroup(IP string, svc *ec2.EC2, securityGroup *ec2.GroupIdentifier) bool {
	paramsSG := &ec2.DescribeSecurityGroupsInput{
		DryRun:   aws.Bool(false),
		GroupIds: []*string{aws.String(*securityGroup.GroupId)},
	}

	resp, err := svc.DescribeSecurityGroups(paramsSG)

	if err != nil {
		panic(err)
	}
	ipExists := false
	for _, ip := range resp.SecurityGroups[0].IpPermissions {
		for _, ipRange := range ip.IpRanges {
			ipExists = strings.Contains(*ipRange.CidrIp, IP)
		}
	}
	return ipExists
}

// Authorize the IP Address on SecurityGroup
func AddIPAddressOnSecurityGroup(svc *ec2.EC2, groupID *string) *ec2.AuthorizeSecurityGroupIngressOutput {
	paramsReq := &ec2.AuthorizeSecurityGroupIngressInput{
		DryRun:     aws.Bool(false),
		GroupId:    aws.String(*groupID),
		CidrIp:     aws.String(IP + "/32"),
		IpProtocol: aws.String("TCP"),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
	}
	out, err := svc.AuthorizeSecurityGroupIngress(paramsReq)
	if err != nil {
		panic(err)
	}
	return out
}
