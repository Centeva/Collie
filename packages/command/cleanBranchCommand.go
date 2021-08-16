package command

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"bitbucket.org/centeva/collie/packages/external"
)

type CleanBranchCommand struct {
	CleanBranch string `tc:"cleanbranch"`
}

func (c *CleanBranchCommand) IsCurrentSubcommand() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "CleanBranch")
}

func (c *CleanBranchCommand) GetFlags(flagProvider external.IFlagProvider) (err error) {
	cmd := flagProvider.NewFlagSet("CleanBranch", "Format branch name Usage: CleanBranch <Branch>")
	cmd.Parse(os.Args[1:])
	if len(os.Args) <= 2 || os.Args[2] == "" {
		return errors.New("format branch name Usage: CleanBranch <Branch>")
	}
	c.CleanBranch = os.Args[2]
	return
}

func (c *CleanBranchCommand) FlagsValid() (err error) {
	if c.CleanBranch == "" {
		return errors.New("CleanBranch is required")
	}
	return
}

func (c *CleanBranchCommand) Execute(globals *GlobalCommandOptions) (err error) {
	name := CleanBranch(c.CleanBranch)
	switch globals.Logger {
	case TEAMCITY:
		{
			paramName, err := GetTeamcityTag(c, "CleanBranch")
			if err != nil {
				return err
			}
			log.Printf("##teamcity[setParameter name='%s' value='%s']", paramName, name)
		}
	case CLI:
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
	matchSlash := regexp.MustCompile(`/`)
	matchSpecial := regexp.MustCompile(`[^\w\s-]`)
	res := matchSlash.ReplaceAllString(name, "-")
	res = matchSpecial.ReplaceAllString(res, "")
	res = strings.ToLower(res)
	return res
}
