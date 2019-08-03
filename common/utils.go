package common

import (
	"fmt"

	"github.com/greenled/portainer-stack-utils/util"

	"github.com/greenled/portainer-stack-utils/client"
)

func GetStackByName(name string) (stack client.Stack, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}

	stacks, err := client.GetStacks("", 0)
	if err != nil {
		return
	}

	util.PrintVerbose(fmt.Sprintf("Getting stack %s...", name))
	for _, stack := range stacks {
		if stack.Name == name {
			return stack, nil
		}
	}
	err = &StackNotFoundError{
		StackName: name,
	}
	return
}

// Custom customerrors
type StackNotFoundError struct {
	StackName string
}

func (e *StackNotFoundError) Error() string {
	return fmt.Sprintf("Stack %s not found", e.StackName)
}
