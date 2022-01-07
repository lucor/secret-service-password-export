package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/godbus/dbus/v5"
	keyring "github.com/ppacher/go-dbus-keyring"
)

func main() {

	log.SetFlags(0)

	var format string
	var name string
	var output string
	var help bool

	flag.StringVar(&name, "c", "", "")
	flag.StringVar(&name, "collection", "", "")
	flag.StringVar(&format, "f", "paw", "")
	flag.StringVar(&format, "format", "paw", "")
	flag.StringVar(&output, "o", "", "")
	flag.StringVar(&output, "output", "", "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	flag.Parse()

	if help {
		log.Println(`Usage:
	secret-service-export [collection]
	
Options:
	-c, --collection	Collection to export. Leave empty to list the available collections
	-f, --format		Output format for the export. Allowed values: [paw, csv]. Default to Paw JSON format
	-o, --output		Write the output to the specified file. If omitted, writes to stdout
	-h, --help		Displays the help and exit

Export the Secrect Service collection using the specified format to stdout.`)
		os.Exit(0)
	}

	bus, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("could not open D-Bus session: %s", err)
	}

	secrets, err := keyring.GetSecretService(bus)
	if err != nil {
		log.Fatalf("could not create the client for the Secret Service: %s", err)
	}

	if name == "" {
		collections, err := secrets.GetAllCollections()
		if err != nil {
			log.Fatalf("could not retrieve the collections: %s", err)
		}

		for _, collection := range collections {
			label, err := collection.GetLabel()
			if err != nil {
				log.Print(err)
			}
			if label == "" {
				continue
			}
			fmt.Println(label)
		}
		os.Exit(0)
	}

	collection, err := secrets.GetCollection(name)
	if err != nil {
		log.Fatal(err)
	}

	items, err := collection.GetAllItems()
	if err != nil {
		log.Fatal(err)
	}

	session, err := secrets.OpenSession()
	if err != nil {
		log.Fatal(err)
	}

	data := []pawLogin{}
	t := time.Now().Format(time.RFC1123)
	for _, item := range items {
		// make sure it is unlocked
		// this also handles any prompt that may be required
		_, err := item.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		label, err := item.GetLabel()
		if err != nil {
			log.Fatal(err)
		}

		secret, err := item.GetSecret(session.Path())
		if err != nil {
			log.Fatal(err)
		}

		if len(secret.Value) == 0 {
			log.Printf("skipped item with empty password: %s", label)
			continue
		}

		created, err := item.GetCreated()
		if err != nil {
			log.Fatal(err)
		}
		modified, err := item.GetCreated()
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, pawLogin{
			Password: pawPassword{
				Value: string(secret.Value),
			},
			Metadata: pawMetadata{
				Name:     label,
				Type:     loginItemType,
				Created:  created,
				Modified: modified,
			},
			Note: pawNote{
				Value: fmt.Sprintf("exported from the Secret Service collection %q\n(%s)", name, t),
			},
		})
	}

	w := os.Stdout
	if output != "" {
		w, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
		if err != nil {
			log.Fatalf("could not create the output file: %s", err)
		}
	}
	defer w.Close()

	if format == "csv" {
		csvw := csv.NewWriter(w)
		csvw.Write([]string{"name", "password", "created", "modified"})
		for _, v := range data {
			csvw.Write([]string{v.Metadata.Name, v.Password.Value, v.Metadata.Created.String(), v.Metadata.Modified.String()})
		}
		csvw.Flush()
		if err := csvw.Error(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// Default to Paw JSON export
	v := map[string][]pawLogin{
		"login": data,
	}

	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Fatal(err)
	}
}

type pawLogin struct {
	Metadata pawMetadata `json:"metadata,omitempty"`
	Note     pawNote     `json:"note,omitempty"`
	Password pawPassword `json:"password,omitempty"`
}

const loginItemType = 8

type pawPassword struct {
	Value string `json:"value,omitempty"`
}

type pawMetadata struct {
	// Title reprents the item name
	Name string `json:"name,omitempty"`
	// Type represents the item type
	Type int `json:"type,omitempty"`
	// Modified holds the modification date
	Modified time.Time `json:"modified,omitempty"`
	// Created holds the creation date
	Created time.Time `json:"created,omitempty"`
}

type pawNote struct {
	Value string `json:"value,omitempty"`
}
