package main

import (
	"fmt"
	"taskWorker/common"
)

func main() {
	env, err := common.GetEnv()
	if nil != err {
		panic(fmt.Errorf("cannot get env variables: %w", err))
	}

	fmt.Printf("%+v\n", env)
}
