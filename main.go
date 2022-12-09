package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

func getRegion(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := imds.NewFromConfig(cfg)
	region, err := client.GetRegion(context.TODO(), &imds.GetRegionInput{})
	if err != nil {
		log.Printf("Unable to retrieve the region from the EC2 instance %v\n", err)
	}

	fmt.Printf("region: %v\n", region)

	res := struct {
		Region string `json:"region"`
	}{}
	res.Region = region.Region
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/", getRegion)
	err := http.ListenAndServe(":80", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("run")
}
