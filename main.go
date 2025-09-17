package main

import "log"

func main() {
	loadEnv()
	connectDB()
	log.Println("🚀 Starting API ...")
	handleRequests()
}
