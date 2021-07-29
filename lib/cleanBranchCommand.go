package lib

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

type CleanBranchCommandOptions struct {
}
type CleanBranchCommand struct {
	CleanBranch string `tc:"cleanbranch"`
}

func (c *CleanBranchCommand) GetFlags(flagProvider IFlagProvider) {
	flagProvider.StringVar(&c.CleanBranch, "CleanBranch", "", "Name of a a branch to format")
}

func (c *CleanBranchCommand) FlagsValid() bool {
	return c.CleanBranch != ""
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
