package awstest

import (
	"fmt"
	"testing"

	"github.com/mgutz/logxi/v1"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestLogin(t *testing.T) {

	// Setup a session configuration that uses the environment variables
	// typically used for AWS values
	//
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	//Load EC2 instance catalog to exercise the basic AWS functionality
	//
	resp, err := ec2.New(sess).DescribeInstances(nil)
	if err != nil {
		log.Warn(fmt.Sprintf("aws failed due to %s", err.Error()), "error", err)
		t.Fatal(err)
	}

	// Print results if asked
	//
	if log.IsDebug() {
		log.Debug(fmt.Sprintf("%v", resp))
	} else {
		instances := 0
		for _, res := range resp.Reservations {
			instances += len(res.Instances)
		}
		log.Info(fmt.Sprintf("%d EC2 instances found", instances))
	}
}
