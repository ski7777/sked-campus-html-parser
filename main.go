package main

import (
	"github.com/ski7777/sked-campus-html-parser/pkg/timetablepage"
	"log"
)

func main() {
	ttp, err := timetablepage.ParseHTMLURL("https://www.asw-ggmbh.de/fileadmin/download/download/Sked%20Stundenplan/Studium/DBWINFO-A02-6.%20Block.html")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("success", len(ttp.TimeTables))
	}
}
