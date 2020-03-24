// Package records interfaces the protobuf definitions with Scribe
package records

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/will-rowe/scribe/src/backend"
)

// error messages
var (
	// ErrNotFound is returned by operations that can't locate the required project
	ErrNotFound = errors.New("project not found")
)

// InitDB will init the project database
func InitDB() *ProjectDatabase {

	// create the db
	db := &ProjectDatabase{
		Projects: make(map[string]*Project),
		Pin:      true,
	}

	// return pointer to the db
	return db
}

// GetNumProjects will return the number of projects in the project database
func (db *ProjectDatabase) GetNumProjects() int {
	return len(db.Projects)
}

// Pull will pull a database from the IPFS using the provided CID
func (db *ProjectDatabase) Pull(node *backend.Node, cid string) error {
	if len(cid) < 1 {
		return fmt.Errorf("no CID provided")
	}

	// prevent over write of db
	if len(db.Projects) != 0 {
		return fmt.Errorf("project database is not empty, Pull method does not currently support merges")
	}

	// get the DAG and load into the struct
	return node.DagGet(cid, "", db)

}

// Push will push the database to the IPFS and return the CID (and any error)
func (db *ProjectDatabase) Push(node *backend.Node) (string, error) {

	// marshal the db as json
	buf := &bytes.Buffer{}
	jsonMarshaller := jsonpb.Marshaler{
		EnumsAsInts:  false, // Whether to render enum values as integers, as opposed to string values.
		EmitDefaults: false, // Whether to render fields with zero values
		Indent:       "\t",  // A string to indent each level by
		OrigName:     false, // Whether to use the original (.proto) name for fields
	}
	if err := jsonMarshaller.Marshal(buf, db); err != nil {
		return "", err
	}

	// add to IPFS
	cid, err := node.DagPut(buf.Bytes(), "json", "cbor", db.Pin)
	if err != nil {
		return "", err
	}

	// return the CID
	return cid, nil
}

// AddProject will add a project to the db
func (db *ProjectDatabase) AddProject(project *Project) error {

	// check the project label is not already in the db
	if _, exists := db.Projects[project.GetLabel()]; exists {
		return fmt.Errorf("project already in the database (label: %s)", project.GetLabel())
	}

	// add it to the db
	db.Projects[project.GetLabel()] = project
	return nil
}

// GetProject will get a project from the db
func (db *ProjectDatabase) GetProject(projectLabel string) (*Project, error) {
	if project, exists := db.Projects[projectLabel]; exists {
		return project, nil
	}
	return nil, ErrNotFound
}
