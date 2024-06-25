package init

import (
	// This package must import all dialects inside sql and must be imported in order to have their init functions run
	_ "github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql/postgres"
	_ "github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql/sqlite"
)
