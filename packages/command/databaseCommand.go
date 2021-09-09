package command

import (
	"log"
	"os"
	"strings"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type DatabaseCommand struct {
	postgresManager external.IPostgresManager
	cmd             external.IFlagSet

	Database         string
	ConnectionString *string
}

func NewDatabaseCommand(flagProvider external.IFlagProvider, postgresManager external.IPostgresManager) *DatabaseCommand {
	return &DatabaseCommand{
		postgresManager: postgresManager,
		cmd:             flagProvider.NewFlagSet("DeleteDatabase", "Delete a Postgres Database Usage: DeleteDatabase <database> [args]"),
	}
}

func (d *DatabaseCommand) IsCurrent() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "DeleteDatabase")
}

func (d *DatabaseCommand) GetFlags() (err error) {
	d.ConnectionString = d.cmd.String("ConnectionString", "", "(required) Postgres database connectionString")

	if len(os.Args) <= 2 || os.Args[2] == "" {
		return errors.New("DeleteDatabase must have a database name")
	}
	d.Database = os.Args[2]

	d.cmd.Parse(os.Args[3:])
	return
}

func (d *DatabaseCommand) Execute() (err error) {
	err = d.postgresManager.Connect(*d.ConnectionString)

	if err != nil {
		return errors.Wrap(err, "Execute failed to connect")
	}
	defer d.postgresManager.Close()
	err = d.postgresManager.DeleteDatabase(d.Database)

	if err != nil {
		return errors.Wrap(err, "Execute failed to delete database")
	}

	log.Printf("Database %s deleted", d.Database)
	return
}
