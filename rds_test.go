package awstest

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/mgutz/logxi/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jackc/pgx"
)

var conn *pgx.Conn

// TestBasicRDS is used to exercise the AWS rds functions for creating, and destroy
// postgres Databases hosted by the AWS rds service
//
func TestBasicRDS(t *testing.T) {

	// use the users name as the AWS DB and resource associated user name so things look somewhat rational if
	// an AWS admin wishes to ask about any left over resources from the test
	osAccount, err := user.Current()
	if err != nil {
		log.Warn(fmt.Sprintf("unable to determine the users local name due to %s", err.Error()), "error", err)
		osAccount = &user.User{Username: pseudoUUID()}
	}

	// Setup a session configuration that uses the environment variables
	// typically used for AWS values
	//
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Now create a security group that exists within the default VPC just for the duration of this test.
	// Also, defer the destroyer for it
	//
	secID, err := newSecGrp(sess, osAccount.Username+" "+pseudoUUID())
	if err != nil {
		t.Fatal(err.Error())
	}
	defer delSecGrp(secID)

	// Use RDS to create a databases with the ID set to a UUID for testing purposes
	dbID := "DB-" + pseudoUUID()
	password := pseudoUUID()
	size := aws.Int64(5)

	params := &rds.CreateDBInstanceInput{
		DBInstanceClass:      aws.String("db.t2.micro"),
		DBInstanceIdentifier: aws.String(dbID),
		Engine:               aws.String("postgres"),
		EngineVersion:        aws.String("9.6.1"),
		StorageType:          aws.String("gp2"),
		AllocatedStorage:     size, // In GB
		AvailabilityZone:     aws.String(*sess.Config.Region + "a"),
		VpcSecurityGroupIds: []*string{
			aws.String(secID),
		},
		MasterUsername:        aws.String(osAccount.Username),
		MasterUserPassword:    aws.String(password),
		BackupRetentionPeriod: aws.Int64(0),
		MultiAZ:               aws.Bool(false),
	}
	resp, err := rds.New(sess).CreateDBInstance(params)
	if err != nil {
		msg := fmt.Sprintf("rds Postgres instance could not be created due to %s", err.Error())
		log.Error(msg, "error", err)
		t.Fatal(msg)
	}
	log.Debug(fmt.Sprintf("%v", resp))

	// Destroy the Test DB that was created
	delResp, err := rds.New(sess).DeleteDBInstance(&rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(dbID),
		SkipFinalSnapshot:    aws.Bool(true),
	})

	if err != nil {
		msg := fmt.Sprintf("rds Postgres instance %s could not be destroyed due to %s", dbID, err.Error())
		log.Error(msg, "error", err)
		t.Fatal(msg)
	}
	log.Debug(fmt.Sprintf("%v", delResp))
}
