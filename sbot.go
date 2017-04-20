package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/ezdiy/go-ssb"
	"github.com/ezdiy/go-ssb/gossip"
	shs "github.com/ezdiy/secretstream/secrethandshake"
	"github.com/andyleap/boltinspect"
)


func main() {
	var kp *shs.EdKeyPair
	var err error
	if len(os.Args)>1 {
		kp, err = shs.LoadSSBKeyPair(os.Args[1])
	}
	if err != nil {
		kp, err = shs.GenEdKeyPair(nil)
	}
	ds, _ := ssb.OpenDataStore("feeds.db", kp)
	fmt.Println("We're ", ds.PrimaryRef)
	gossip.Replicate(ds,"")

	bi := boltinspect.New(ds.DB())
	http.HandleFunc("/", bi.InspectEndpoint)
	http.ListenAndServe(":45079", nil)
}


