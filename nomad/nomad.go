package nomad

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/hashicorp/nomad/api"
)

type Nomad struct {
	Address    string
	AllocsDir  string
	NodeID     string
	MetaPrefix string
}

func (n *Nomad) Client() *api.Client {
	config := *api.DefaultConfig()
	config.Address = n.Address
	client, err := api.NewClient(&config)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (n *Nomad) SetNodeIDFromEnvs() {
	q := &api.QueryOptions{Namespace: os.Getenv("NOMAD_NAMESPACE")}
	alloc, _, err := n.Client().Allocations().Info(os.Getenv("NOMAD_ALLOC_ID"), q)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Found node id %s using env vars\n", alloc.NodeID)
	n.NodeID = alloc.NodeID
}

func (n *Nomad) Allocs() []*api.Allocation {
	allocs, _, err := n.Client().Nodes().Allocations(n.NodeID, nil)
	if err != nil {
		panic(err)
	}
	return allocs
}

func (n *Nomad) TaskMeta(Task api.Task) map[string]string {
	meta := make(map[string]string)

	regex, _ := regexp.Compile(fmt.Sprintf("^(%s)\\.", n.MetaPrefix))
	for key, value := range Task.Meta {
		if regex.MatchString(key) {
			strippedKey := regex.ReplaceAllString(key, "")
			meta[strippedKey] = value
		}
	}

	return meta
}

func (n *Nomad) TaskMetaGet(Task api.Task, Key string, Default string) string {
	meta := n.TaskMeta(Task)

	value, exists := meta[Key]

	if exists {
		return value
	}

	return Default
}

func AllocTasks(Alloc *api.Allocation) ([]*api.Task, error) {
	for _, group := range Alloc.Job.TaskGroups {
		if *group.Name == Alloc.TaskGroup {
			return group.Tasks, nil
		}
	}

	return nil, errors.New("could not find Tasks for Allocation")
}
