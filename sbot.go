package main

import (
	"os"
	"fmt"
	"bytes"
	"net/http"
	"crypto/sha256"
	"github.com/ezdiy/go-ssb"
	"github.com/ezdiy/go-ssb/gossip"
	shs "github.com/ezdiy/secretstream/secrethandshake"
	"github.com/andyleap/boltinspect"
	_ "net/http/pprof"
)


func main() {
	var kp *shs.EdKeyPair
	var name = os.Args[1]
	kp, _ = shs.LoadSSBKeyPair(name)
	if kp == nil {
		h := sha256.Sum256([]byte(name))
		kp, _ = shs.GenEdKeyPair(bytes.NewReader(h[:]))
	}
	ds, _ := ssb.OpenDataStore(nil, name + ".db", kp, 0)
	me := ds.GetFeed(ds.PrimaryRef)
	for _, m := range os.Args[2:] {
		fmt.Println("publishing ", m)
		me.PublishMessageJSON([]byte(m))
	}
	fmt.Println("We're ", ds.PrimaryRef)
	gossip.Replicate(ds,"")

	bi := boltinspect.New(ds.DB())
	http.HandleFunc("/bolt", bi.InspectEndpoint)
	http.ListenAndServe(":45079", nil)
}


