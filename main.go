package main

import (
	"fmt"

	"github.com/gdexlab/go-render/render"
	stores "github.com/hxhieu/oceanbridge-wrapper-go-common/persistent/userstores"
)

func main() {
	// Example usage of the stores

	userStore, err := stores.NewUserStoreFirestore()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = userStore.SetToken("56balI5EXZV6JRPV7ZPe", "hello tokenasdas")

	if err != nil {
		fmt.Println(render.AsCode(err))
	}
}
