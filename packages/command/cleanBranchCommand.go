package command

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type CleanBranchCommand struct {
	cmd external.IFlagSet

	CleanBranch string `tc:"cleanbranch"`
	Logger      *string
}

func NewCleanBranchCommand(flagProvider external.IFlagProvider) *CleanBranchCommand {
	return &CleanBranchCommand{
		cmd: flagProvider.NewFlagSet("CleanBranch", "Format branch name Usage: CleanBranch <Branch>"),
	}
}

func (c *CleanBranchCommand) IsCurrent() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "CleanBranch")
}

func (c *CleanBranchCommand) GetFlags() (err error) {
	c.Logger = c.cmd.String("Logger", string(CLI), "Log output style to use [cli|teamcity]")

	if len(os.Args) <= 2 || os.Args[2] == "" {
		c.cmd.PrintDefaults()
		return errors.New("Missing branch, see usage.")
	}

	c.CleanBranch = os.Args[2]

	c.cmd.Parse(os.Args[3:])

	if c.CleanBranch == "" {
		return errors.New("CleanBranch is required")
	}

	if c.Logger == nil || *c.Logger == "" {
		return errors.New("logger must have a value")
	}

	logger := *c.Logger

	if logger != "cli" && logger != "teamcity" {
		return errors.Errorf("logger must be either 'cli' or 'teamcity' got '%s'", logger)
	}
	return
}

func (c *CleanBranchCommand) Execute() (err error) {
	name := CleanBranch(c.CleanBranch)
	logger := *c.Logger
	switch logger {
	case string(TEAMCITY):
		{
			paramName, err := GetTeamcityTag(c, "CleanBranch")
			if err != nil {
				return err
			}
			log.Printf("##teamcity[setParameter name='%s' value='%s']", paramName, name)
		}
	case string(CLI):
		fallthrough
	default:
		log.Printf("%s", name)
	}

	return
}

func GetTeamcityTag(kind interface{}, fieldName string) (paramName string, err error) {

	field, ok := reflect.TypeOf(kind).Elem().FieldByName(fieldName)

	if !ok {
		return "", fmt.Errorf("could not find Field %s on Type %v", fieldName, kind)
	}
	paramName, ok = field.Tag.Lookup("tc")

	if !ok {
		return "", fmt.Errorf("field: %s does not have a tc tag", fieldName)
	}

	return
}

func CleanBranch(name string) string {
	matchSlash := regexp.MustCompile(`[/_]`)
	matchSpecial := regexp.MustCompile(`[^\w\s-]`)
	res := matchSlash.ReplaceAllString(name, "-")
	res = matchSpecial.ReplaceAllString(res, "")
	res = strings.ToLower(res)
	return res
}
