package russifiers

import (
	"fmt"
	"strings"
)

type JudokaName struct {
	FirstName    string
	LastName     string
	FirstNameRus *string
	LastNameRus  *string
}

func NewJudokaName(firstName, lastName string, firstNameRus, lastNameRus *string) JudokaName {
	return JudokaName{
		FirstName:    firstName,
		LastName:     lastName,
		FirstNameRus: firstNameRus,
		LastNameRus:  lastNameRus,
	}
}

type JudokaRussifier map[string][]string

func NewJudokaRussifier(judokas []JudokaName) JudokaRussifier {
	russifier := make(map[string][]string)

	for _, judoka := range judokas {
		fullNameEng := fmt.Sprintf("%s %s", judoka.LastName, judoka.FirstName)
		if judoka.LastNameRus != nil && judoka.FirstNameRus != nil {
			russifier[fullNameEng] = []string{*judoka.FirstNameRus, *judoka.LastNameRus}
		} else {
			russifier[fullNameEng] = []string{judoka.FirstName, judoka.LastName}
		}
	}
	return JudokaRussifier(russifier)
}

func (j JudokaRussifier) Russify(fullName string) []string {
	if russified, ok := j[fullName]; ok {
		return russified
	}

	return strings.Split(fullName, " ")
}
