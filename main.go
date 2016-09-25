package main

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	elastic "gopkg.in/olivere/elastic.v3"
)

const mapping = `{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "mappings":{
        "samples":{
            "properties":{
                "scan_date": {
                  "type": "date"
                },
                "file":{
                    "type":"string"
                },
                "plugins":{
                    "type":"string"
                }
            }
        }
    }
}`

func main() {
	client, err := elastic.NewSimpleClient()
	utils.Assert(err)

	exists, err := client.IndexExists("malice").Do()
	utils.Assert(err)

	if !exists {
		// Index does not exist yet.
		createIndex, err := client.CreateIndex("malice").BodyString(mapping).Do()
		utils.Assert(err)
		if !createIndex.Acknowledged {
			// Not acknowledged
			log.Error("Couldn't create Index.")
		} else {
			log.Info("Created Index: ", "malice")
		}
	} else {
		log.Info("Index malice already exists.")
	}

	// sample, err := client.Get().
	// 	Index("malice").
	// 	Type("samples").
	// 	Id("1").
	// 	Do()
	//
	// fmt.Println(sample)
	// fmt.Println(err)
	// if err != nil {
	//
	// }

	// if sample.Found {
	// 	fmt.Printf("Got document %s in version %d from index %s, type %s\n", sample.Id, sample.Version, sample.Index, sample.Type)
	// } else {

	scan := map[string]interface{}{
		// "id":      sample.SHA256,
		"file":      "file",
		"plugins":   "plugins",
		"scan_date": time.Now().Format(time.RFC3339Nano),
	}
	jsonString, err := json.Marshal(scan)
	utils.Assert(err)

	log.Infoln(string(jsonString))

	newSample, err := client.Index().
		Index("malice").
		Type("samples").
		OpType("create").
		// Id("1").
		BodyJson(scan).
		Do()
	utils.Assert(err)
	log.Infof("Indexed sample %s to index %s, type %s\n", newSample.Id, newSample.Index, newSample.Type)

	update, err := client.Update().Index("malice").Type("samples").Id(newSample.Id).
		Doc(map[string]interface{}{"plugins": "suck"}).
		Do()
	utils.Assert(err)
	fmt.Printf("New version of sample %q is now %d\n", update.Id, update.Version)

	// }

}