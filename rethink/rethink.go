package rethink

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx/tracing"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const tagName = "db"

// Config represents rethink database config.
type Config struct {
	Addr string
	Name string
	User string
	Pass string
}

// Client represents rethink database client.
type Client struct {
	db     *r.Session
	tables []string

	tracing *tracing.Tracer
}

// New create new rethink database instance.
func New(c Config) (*Client, error) {
	sess, err := r.Connect(r.ConnectOpts{
		Address:    c.Addr,
		Database:   c.Name,
		Timeout:    time.Minute,
		InitialCap: 10,
		MaxOpen:    10,
	})
	if err != nil {
		return nil, err
	}

	r.SetTags(tagName)

	ts, err := getTables(sess)
	if err != nil {
		log.Fatal("could not get table:", err)
	}

	return &Client{sess, ts, nil}, nil
}

func (c *Client) SetTracer(t *tracing.Tracer) {
	c.tracing = t
}

// Close ends rethink database session.
func (c *Client) Close() error {
	return c.db.Close()
}

// autoMigrate create tables and wait for to finish. Existing table will be ignored.
func (c *Client) autoMigrate(table string) error {
	// Checks table existence and skip create.
	for _, t := range c.tables {
		if t == table {
			return nil
		}
	}

	return c.exec(r.TableCreate(table))
}

// autoIndex creates table index base model that has tag "index".
func (c *Client) autoIndex(table string, model interface{}) error {
	for _, ff := range getModelIndexedFields(model) {
		fmt.Println("autoIndex", table, ff)
		if err := c.createIndex(table, ff); err != nil {
			return fmt.Errorf("could not create %s index on %s table: %s", ff, tableCatalog, err)
		}
	}

	return nil
}

// run returns a cursor which can be used to view all rows returned.
func (c *Client) run(t r.Term) (*r.Cursor, error) {
	return t.Run(c.db)
}

// runWrite returns a WriteResponse and should be used for queries such as Insert, Update, etc...
func (c *Client) runWrite(t r.Term) (r.WriteResponse, error) {
	return t.RunWrite(c.db)
}

// exec sends a query to the server and closes the connection immediately after reading the response from the database.
// If you do not wish to wait for the response then you can set the NoReply flag.
func (c *Client) exec(t r.Term) error {
	return t.Exec(c.db)
}

func (c *Client) list(t r.Term, out interface{}) error {
	if c.tracing != nil {
		s := c.tracing.StartSpan("rethink list " + t.String())
		defer s.End()
	}

	res, err := c.run(t)
	if err != nil {
		return err
	}
	if err = res.All(out); err != nil {
		return err
	}

	return res.Close()
}

func (c *Client) one(t r.Term, out interface{}) error {
	if c.tracing != nil {
		s := c.tracing.StartSpan("rethink one " + t.String())
		defer s.End()
	}

	res, err := c.run(t)
	if err != nil {
		return err
	}
	if err = res.One(out); err != nil {
		return err
	}

	return res.Close()
}

func (c *Client) insert(t r.Term) (id string, err error) {
	res, err := c.runWrite(t)
	if err != nil {
		return
	}

	if len(res.GeneratedKeys) == 0 {
		return "", nil
	}

	return res.GeneratedKeys[0], nil
}

func (c *Client) update(t r.Term) error {
	_, err := c.runWrite(t)
	return err
}

func (c *Client) delete(t r.Term) error {
	return c.update(t)
}

func (c *Client) createIndex(tableName, index string) error {
	tbl := r.Table(tableName)

	var indexes []string
	if err := c.list(tbl.IndexList(), &indexes); err != nil {
		return err
	}

	for _, ii := range indexes {
		// Skip creating index.
		if ii == index {
			return nil
		}
	}

	return c.exec(tbl.IndexCreate(index))
}

func getTables(s *r.Session) (table []string, err error) {
	res, _ := r.TableList().Run(s)
	err = res.All(&table)
	return
}

func now() *time.Time {
	t := time.Now()
	return &t
}

const indexedTagName = "index"

func getModelIndexedFields(model interface{}) (fields []string) {
	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(model)
	// Iterate over all available fields and read the tag value.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if !strings.Contains(tag, indexedTagName) {
			continue
		}

		tagField := strings.Split(tag, ",")[0]
		// Ignore ID field since its index by default.
		if tagField == "id" {
			continue
		}

		// Only get the base field name.
		fields = append(fields, tagField)
	}

	return
}
