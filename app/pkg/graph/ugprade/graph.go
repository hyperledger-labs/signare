// Package upgrade defines the graph with the dependencies needed to execute the db migrations.
package upgrade

import (
	"github.com/asaskevich/govalidator"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql"
)

type UpgradeGraph struct {
	config Config

	librariesGraph *librariesGraph
}

func New(config Config) (*UpgradeGraph, error) {
	valid, err := govalidator.ValidateStruct(config)
	if err != nil || !valid {
		return nil, err
	}

	return &UpgradeGraph{
		config: config,
	}, nil
}

// Build builds the upgrade graph
func (graph *UpgradeGraph) Build() {
	var err error

	// Libraries
	graph.librariesGraph, err = initializeLibraries(graph.config)
	checkError(err)
}

func (graph *UpgradeGraph) PersistenceFwConnection() sql.Connection {
	return graph.librariesGraph.persistenceConnection
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
