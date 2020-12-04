package main

import (
	"log"
	"os"
	"os/signal"
)

type Code int

func main() {
	KillSignal := make(chan os.Signal, 1)
	signal.Notify(KillSignal, os.Interrupt, os.Kill)

	cfg = new(Config)
	cfg.Load()

	CyT = new(CycleTask)
	CyT.Load()

	http := new(HTTPServer)
	http.StartHTTPServer()

	<-KillSignal
	CyT.Stop()
	CyT.Save()
	http.srv.Close()
	// http.srv.Shutdown(context.TODO())
	http.Wait()
	cfg.Save()
	log.Println("Thanks!")
}
