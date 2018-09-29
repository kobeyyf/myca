package main

import (
	"errors"
	"fmt"
)

type Action interface {
	Check(args *Args) error
	Run() error
}

func main() {
	args, err := GetArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	action, err := newAction(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := action.Check(args); err != nil {
		fmt.Println(err)
		return
	}
	if err := action.Run(); err != nil {
		fmt.Println(err)
		// fmt.Println(args)
		return
	}
}

func newAction(args *Args) (Action, error) {
	switch args.command {
	case "Init":
		return new(ActionInit), nil
	case "AddUser":
		return new(ActionAddUser), nil
	case "AddOrderer":
		return new(ActionAddOrderer), nil
	case "AddPeer":
		return new(ActionAddPeer), nil
	case "AddOrg":
		return new(ActionAddOrg), nil
	}
	return nil, errors.New("unSupportCommand")
}
