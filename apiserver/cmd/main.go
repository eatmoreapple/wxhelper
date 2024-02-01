package main

import "github.com/eatmoreapple/wxhelper/apiserver"

const defaultADDR = ":19089"

func main() {
	// todo add flag here
	srv := apiserver.Default()
	srv.Run(defaultADDR)
}
