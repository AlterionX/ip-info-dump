package infosource

import (
	"flag"
	"fmt"
	"log"
	"os"

	api "github.com/VirusTotal/vt-go"
)

var apikey = flag.String("apikey", "", "VirusTotal API key")
var sha256 = flag.String("sha256", "", "SHA-256 of some file")

func basicVT() {

	flag.Parse()

	if *apikey == "" || *sha256 == "" {
		fmt.Println("Must pass both the --apikey and --sha256 arguments.")
		os.Exit(0)
	}

	client := api.NewClient(*apikey)

	file, err := client.GetObject(api.URL("files/%s", *sha256))
	if err != nil {
		log.Fatal(err)
	}

	ls, err := file.GetTime("last_submission_date")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File %s was submitted for the last time on %v\n", file.ID(), ls)
}

type virustotal struct{}
