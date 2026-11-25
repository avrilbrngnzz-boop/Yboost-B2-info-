package main

import yboost "yboost/src"

func main() {
	yboost.InitLogger()
	yboost.CreateDB() //a faire
	yboost.StartServer()
}
