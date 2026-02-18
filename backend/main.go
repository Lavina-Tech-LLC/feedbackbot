package main

import (
	_ "github.com/Lavina-Tech-LLC/feedbackbot/internal/db"
	webServer "github.com/Lavina-Tech-LLC/feedbackbot/internal/webserver"
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func main() {
	go webServer.Listen()

	lvn.WaitExitSignal()
}
