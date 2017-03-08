package awstest

// The file contains functions used for prototyping postgres rds
// functions.  More information about aws rds postgres support can be found
// at https://aws.amazon.com/rds/postgresql/
//
import (
	"fmt"
	"time"

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
func CreateDB(ec2Type string, allocation int64, userName string, password string) (dbID string, descr *rds.DBInstance, cleanup func(snapshot bool) (err error), err error) {

	dbID = "DB-" + PseudoUUID()

	if dbID, cleanup, err = createRDS(dbID, userName, password, ec2Type, allocation); err != nil {
		return dbID, nil, cleanup, err
	}

	// Setup a session configuration that uses the environment variables
	// typically used for AWS values
	//
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	for {
		descr, err := rds.New(sess).DescribeDBInstances(&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(dbID),
			MaxRecords:           aws.Int64(20),
		})
		if err != nil {
			return dbID, nil, cleanup, err
		}

		dbInst := &rds.DBInstance{}

		for _, aDB := range descr.DBInstances {
			if dbID == *aDB.DBInstanceIdentifier {
				dbInst = aDB
				break
			}
		}

		if len(*dbInst.DBInstanceIdentifier) != 1 {
			return dbID, nil, cleanup, fmt.Errorf("the AWS RDS Instance %s successfully began creation but could never be accessed", dbID)
		}

		if *dbInst.DBInstanceStatus != "creating" {
			return dbID, dbInst, cleanup, err
		}
		select {
		case <-time.After(3 * time.Second):
		}
	}
}

// createRDS is used internally to generate the RDS hosting for the postgres DB
//
// The cleanup function should be called if the caller wishes to destroy the DB
//
func createRDS(instanceID string, userName string, password string, ec2Type string, allocation int64) (
	rdsID string, cleanup func(snapshot bool) (err error), err error) {

	cleanup = func(snapshot bool) (err error) { return nil }

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

	cleanup = func(snapshot bool) (err error) {
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

		// Wait until AWS beings the deletion action which frees up us to remove the
		// dependent security group
		for {
			descr, err := rds.New(sess).DescribeDBInstances(&rds.DescribeDBInstancesInput{
				DBInstanceIdentifier: aws.String(instanceID),
				MaxRecords:           aws.Int64(20),
			})
			if err != nil {
				break
			}

			if len(descr.DBInstances) != 1 {
				break
			}

			if *descr.DBInstances[0].DBInstanceStatus == "deleting" {
				break
			}
			select {
			case <-time.After(3 * time.Second):
			}
		}
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

	_, err = rds.New(sess).CreateDBInstance(params)

	// there is nothing we can sensibly do at this moment to log etc the error, have the caller do it
	return rdsID, cleanup, err
}
