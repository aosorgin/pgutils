package main

import "gitlab.com/aosorgin/pggen/cmd/pgquerycmd"

func main() {
	if err := pgquerycmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
