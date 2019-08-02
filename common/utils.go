package common

import (
	"fmt"
)

func GetStackByName(name string) (stack Stack, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}

	stacks, err := client.GetStacks("", 0)
	if err != nil {
		return
	}

	PrintVerbose(fmt.Sprintf("Getting stack %s...", name))
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

type StackListFilter struct {
	SwarmId    string `json:",omitempty"`
	EndpointId uint32 `json:",omitempty"`
}

// Custom customerrors
type StackNotFoundError struct {
	StackName string
}

func (e *StackNotFoundError) Error() string {
	return fmt.Sprintf("Stack %s not found", e.StackName)
}
