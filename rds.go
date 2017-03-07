package awstest

// The file contains functions used for prototyping postgres rds
// functions.  More information about aws rds postgres support can be found
// at https://aws.amazon.com/rds/postgresql/
//
import (
	"os/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

// CreateDB can be used to create a DB hosted ln AWS rds.
//
// ec2Type must be one of the supported types for rds, further information can
// be found at http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html.
// The ec2 instance used for the compute can be modified usibng the AWS console as a single
// push button action.
//
// The allocation value is the number of Gigabytes that will be reserved for the DB.  Sizing for DB
// is done upfront.  Increasing the sizing can be easily done using push button upgrades to rds
//
// The cleanup function if used will destroy all resources associated with the
// database if the applications wishes to do so once it completes its testing.
//
func CreateDB(ec2Type string, allocation int64, dbName string) (dbID string, cleanup func() (err error), err error) {
	// use the users name as the AWS DB and resource associated user name so things look somewhat rational if
	// an AWS admin wishes to ask about any left over resources from the test
	osAccount, err := user.Current()
	if err != nil {
		osAccount = &user.User{Username: pseudoUUID()}
	}

	dbID = "DB-" + pseudoUUID()
	password := pseudoUUID()

	return createRDS(dbID, osAccount.Username, password, ec2Type, allocation)
}

// createRDS is used internally to generate the RDS hosting for the postgres DB
//
// The cleanup function should be called if the caller wishes to destroy the DB
//
func createRDS(instanceID string, userName string, password string, ec2Type string, allocation int64) (rdsID string, cleanup func() (err error), err error) {

	cleanup = func() (err error) { return nil }

	// Setup a session configuration that uses the environment variables
	// typically used for AWS values
	//
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Now create a security group that exists within the default VPC just for the duration of this test.
	// Also, defer the destroyer for it
	//
	secID, err := newSecGrp(sess, userName)
	if err != nil {
		return rdsID, cleanup, err
	}

	cleanup = func() (err error) {
		// Setup a session configuration that uses the environment variables
		// typically used for AWS values
		//
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		// Destroy the Test DB that was created.  Return codes are ignored as there is
		// nothing meaningful we can do with them
		//
		_, err = rds.New(sess).DeleteDBInstance(&rds.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(instanceID),
			SkipFinalSnapshot:    aws.Bool(true),
		})
		delSecGrp(secID)
		return err
	}

	// Use RDS to create a databases with the ID set to a UUID for testing purposes
	size := aws.Int64(allocation)

	params := &rds.CreateDBInstanceInput{
		DBInstanceClass:      aws.String(ec2Type),
		DBInstanceIdentifier: aws.String(instanceID),
		Engine:               aws.String("postgres"),
		EngineVersion:        aws.String("9.6.1"),
		StorageType:          aws.String("gp2"),
		AllocatedStorage:     size, // In GB
		AvailabilityZone:     aws.String(*sess.Config.Region + "a"),
		VpcSecurityGroupIds: []*string{
			aws.String(secID),
		},
		MasterUsername:        aws.String(userName),
		MasterUserPassword:    aws.String(password),
		BackupRetentionPeriod: aws.Int64(0),
		MultiAZ:               aws.Bool(false),
	}

	if _, err = rds.New(sess).CreateDBInstance(params); err != nil {
		return rdsID, cleanup, err
	}

	return instanceID, cleanup, nil
}
