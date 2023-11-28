package main

var httpPort string
var httpsPort string
var domainList map[string]Domain
var domainMain string

func initVariable() {
	httpPort = "6969"
	httpsPort = "7070"
	domainMain = "forgeforce.org"
	domainList = make(map[string]Domain)
}
