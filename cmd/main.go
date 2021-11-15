package main

import "github.com/aosorgin/pgutils/cmd/pgquerycmd"

func main() {
	if err := pgquerycmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
