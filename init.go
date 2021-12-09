package goapp

func init() {
	LoadEnv()
	InitLogger()
	InitAuth()
	InitTracing()
	InitRouter()
}
