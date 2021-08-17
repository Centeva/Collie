package command

import (
	"log"
	"os"
	"strings"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type DatabaseCommand struct {
	postgresManager  external.IPostgresManager
	database         string
	connectionString *string
}

func NewDatabaseCommand(postgresManager external.IPostgresManager) *DatabaseCommand {
	return &DatabaseCommand{
		postgresManager: postgresManager,
	}
}

func (d *DatabaseCommand) GetFlags(FlagProvider external.IFlagProvider) (err error) {
	cmd := FlagProvider.NewFlagSet("DeleteDatabase", "Delete a Postgres Database Usage: DeleteDatabase <database> [args]")
	if len(os.Args) <= 2 || os.Args[2] == "" {
		return errors.New("DeleteDatabase must have a database name")
	}
	d.database = os.Args[2]

	d.connectionString = cmd.String("ConnectionString", "", "(required) Postgres database connectionString")

	cmd.Parse(os.Args[3:])
	return
}

func (d *DatabaseCommand) IsCurrentSubcommand() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "DeleteDatabase")
}

func (d *DatabaseCommand) FlagsValid() (err error) {
	return
}

func (d *DatabaseCommand) Execute(globals *GlobalCommandOptions) (err error) {
	err = d.postgresManager.Connect(*d.connectionString)

	if err != nil {
		return errors.Wrap(err, "Execute failed to connect")
	}
	defer d.postgresManager.Close()
	err = d.postgresManager.DeleteDatabase(d.database)

	if err != nil {
		return errors.Wrap(err, "Execute failed to delete database")
	}

	log.Printf("Database %s deleted", d.database)
	return
}
