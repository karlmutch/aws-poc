package awstest

import (
	"crypto/rand"
	"fmt"

	"github.com/mgutz/logxi/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// PseudoUUID can be used to generate IDs that will be unique enough for testing using
// place holder values that need to be unique for testing APIs that require values
// for keys and names etc
//
func PseudoUUID() (uuid string) {

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Error(fmt.Sprintf("Error: %s", err.Error()), "error", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// getVPC is used to obtain the default VPC for the account being used
//
func getVPC(sess *session.Session) (id string, err error) {

	resp, err := ec2.New(sess).DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return "", err
	}

	for _, vpc := range resp.Vpcs {
		if *vpc.IsDefault {
			return *vpc.VpcId, nil
		}
	}
	return "", err
}

func newSecGrp(sess *session.Session, userName string) (id string, err error) {

	// Locate the default VPC with which the database deployment will be associated
	//
	vpcID, err := getVPC(sess)
	if err != nil {
		msg := fmt.Sprintf("Failed to get VPC IDs due to %s", err)
		log.Warn(msg, "error", err)
		return "", err
	}

	resp, err := ec2.New(sess).CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description: aws.String(fmt.Sprintf("Test group for %s", userName)),
		GroupName:   aws.String(fmt.Sprintf("Transitory Testing Group (%s)", userName)),
		VpcId:       aws.String(vpcID),
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to get VPC IDs due to %s", err)
		log.Warn(msg, "error", err)
		return "", err
	}

	// By AWS AWS Opens up all egress ports but no ingress ports so we should do this just for the duration
	// of the test
	//
	_, err = ec2.New(sess).AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String("0.0.0.0/0"),
		FromPort:   aws.Int64(1),
		ToPort:     aws.Int64(65535),
		GroupId:    aws.String(*resp.GroupId),
		IpProtocol: aws.String("-1"),
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to open the security group for remote access due to %s", err)
		log.Warn(msg, "error", err)
		return "", err
	}

	log.Debug(fmt.Sprintf("%v", resp))

	return *resp.GroupId, nil
}

func delSecGrp(id string) {
	// Setup a session configuration that uses the environment variables
	// typically used for AWS values
	//
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Delete the supplied group ID we can only log errors which we can do nothing about as this is a cleanup function
	//
	if _, err := ec2.New(sess).DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(id),
	}); err != nil {
		log.Warn(fmt.Sprintf("cleaning up security group %s failed due to %s", id, err.Error()), "error", err)
	}
}
