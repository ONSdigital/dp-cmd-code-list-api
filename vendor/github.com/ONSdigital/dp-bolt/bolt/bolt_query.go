package bolt

import (
	"github.com/pkg/errors"
	neo4j "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
)

//go:generate moq -out mock/bolt.go -pkg mock . DBPool NeoConn NeoRows NeoResult

var NonUniqueResult = errors.New("unique result expected but was not")

var ErrNoResults = errors.New("no results were found")

type Result struct {
	Data  []interface{}
	Meta  map[string]interface{}
	Index int
}

type ResultMapper func(r *Result) error

type NeoConn neo4j.Conn
type NeoRows neo4j.Rows
type NeoResult neo4j.Result

// DBPool contains the methods to control access to the Neo4J
// database pool
type DBPool interface {
	OpenPool() (neo4j.Conn, error)
	Close() error
}

type DB struct {
	pool DBPool
}

//New create a new bolt.DB struct.
func New(pool DBPool) *DB {
	return &DB{pool: pool}
}

//Close attempts to close the db connection pool.
func (d *DB) Close() error {
	return d.pool.Close()
}

//QueryForResults executes the provided query to return 1 or more results.
func (d *DB) QueryForResults(query string, params map[string]interface{}, mapResult ResultMapper) error {
	return d.query(query, params, mapResult, false)
}

//QueryForResults executes the provided query to return a single result.
func (d *DB) QueryForResult(query string, params map[string]interface{}, mapResult ResultMapper) error {
	return d.query(query, params, mapResult, true)
}

func (d *DB) query(cypherQuery string, params map[string]interface{}, mapResult ResultMapper, singleResult bool) error {
	conn, err := d.pool.OpenPool()
	if err != nil {
		return errors.WithMessage(err, "error opening neo4j connection")
	}
	defer conn.Close()

	rows, err := conn.QueryNeo(cypherQuery, params)
	if err != nil {
		return errors.WithMessage(err, "error executing neo4j query")
	}
	defer rows.Close()

	index := 0
	numOfResults := 0
results:
	for {
		data, meta, nextNeoErr := rows.NextNeo()
		if nextNeoErr != nil {
			if nextNeoErr == io.EOF {
				break results
			} else {
				return errors.WithMessage(nextNeoErr, "extractResults: rows.NextNeo() return unexpected error")
			}
		}
		numOfResults++
		if singleResult && index > 0 {
			return NonUniqueResult
		}

		if mapResult != nil {
			if err := mapResult(&Result{Data: data, Meta: meta, Index: index}); err != nil {
				return errors.WithMessage(err, "mapResult returned an error")
			}
		}
		index++
	}

	if numOfResults == 0 {
		return ErrNoResults
	}
	return nil
}
