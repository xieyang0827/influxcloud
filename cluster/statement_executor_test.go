package cluster_test

import (
	"io"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/influxdata/influxdb/influxql"
	"github.com/influxdata/influxdb/models"
	// "github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxdb/coordinator"
	"github.com/influxdata/influxdb/tsdb"
)

const (
	// DefaultDatabase is the default database name used in tests.
	DefaultDatabase = "db0"

	// DefaultRetentionPolicy is the default retention policy name used in tests.
	DefaultRetentionPolicy = "rp0"
)

// // Ensure query executor can execute a simple SELECT statement.
// func TestQueryExecutor_ExecuteQuery_SelectStatement(t *testing.T) {
// 	e := DefaultQueryExecutor()

// 	// The meta client should return a single shard owned by the local node.
// 	e.MetaClient.ShardsByTimeRangeFn = func(sources influxql.Sources, tmin, tmax time.Time) (a []meta.ShardInfo, err error) {
// 		return []meta.ShardInfo{{ID: 100, Owners: []meta.ShardOwner{{NodeID: 0}}}}, nil
// 	}

// 	// The TSDB store should return an IteratorCreator for shard.
// 	// This IteratorCreator returns a single iterator with "value" in the aux fields.
// 	e.TSDBStore.ShardIteratorCreatorFn = func(id uint64) influxql.IteratorCreator {
// 		if id != 100 {
// 			t.Fatalf("unexpected shard id: %d", id)
// 		}

// 		var ic IteratorCreator
// 		ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
// 			return &FloatIterator{Points: []influxql.FloatPoint{
// 				{Name: "cpu", Time: int64(0 * time.Second), Aux: []interface{}{float64(100)}},
// 				{Name: "cpu", Time: int64(1 * time.Second), Aux: []interface{}{float64(200)}},
// 			}}, nil
// 		}
// 		ic.FieldDimensionsFn = func(sources influxql.Sources) (fields, dimensions map[string]struct{}, err error) {
// 			return map[string]struct{}{"value": struct{}{}}, nil, nil
// 		}
// 		ic.SeriesKeysFn = func(opt influxql.IteratorOptions) (influxql.SeriesList, error) {
// 			return influxql.SeriesList{
// 				{Name: "cpu", Aux: []influxql.DataType{influxql.Float}},
// 			}, nil
// 		}
// 		return &ic
// 	}

// 	// Verify all results from the query.
// 	if a := ReadAllResults(e.ExecuteQuery(`SELECT * FROM cpu`, "db0", 0)); !reflect.DeepEqual(a, []*influxql.Result{
// 		{
// 			StatementID: 0,
// 			Series: []*models.Row{{
// 				Name:    "cpu",
// 				Columns: []string{"time", "value"},
// 				Values: [][]interface{}{
// 					{time.Unix(0, 0).UTC(), float64(100)},
// 					{time.Unix(1, 0).UTC(), float64(200)},
// 				},
// 			}},
// 		},
// 	}) {
// 		t.Fatalf("unexpected results: %s", spew.Sdump(a))
// 	}
// }

// // Ensure query executor can enforce a maximum series selection count.
// func TestQueryExecutor_ExecuteQuery_MaxSelectSeriesN(t *testing.T) {
// 	e := DefaultQueryExecutor()
// 	e.StatementExecutor.MaxSelectSeriesN = 3

// 	// The meta client should return a two shards on the local node.
// 	e.MetaClient.ShardsByTimeRangeFn = func(sources influxql.Sources, tmin, tmax time.Time) (a []meta.ShardInfo, err error) {
// 		return []meta.ShardInfo{
// 			{ID: 100, Owners: []meta.ShardOwner{{NodeID: 0}}},
// 			{ID: 101, Owners: []meta.ShardOwner{{NodeID: 0}}},
// 		}, nil
// 	}

// 	// This iterator creator returns an iterator that operates on 2 series.
// 	// Reuse this iterator for both shards. This brings the total series count to 4.
// 	var ic IteratorCreator
// 	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
// 		return &FloatIterator{
// 			Points: []influxql.FloatPoint{{Name: "cpu", Time: int64(0 * time.Second), Aux: []interface{}{float64(100)}}},
// 			stats:  influxql.IteratorStats{SeriesN: 2},
// 		}, nil
// 	}
// 	ic.FieldDimensionsFn = func(sources influxql.Sources) (fields, dimensions map[string]struct{}, err error) {
// 		return map[string]struct{}{"value": struct{}{}}, nil, nil
// 	}
// 	ic.SeriesKeysFn = func(opt influxql.IteratorOptions) (influxql.SeriesList, error) {
// 		return influxql.SeriesList{
// 			{Name: "cpu", Aux: []influxql.DataType{influxql.Float}},
// 		}, nil
// 	}
// 	e.TSDBStore.ShardIteratorCreatorFn = func(id uint64) influxql.IteratorCreator { return &ic }

// 	// Verify all results from the query.
// 	if a := ReadAllResults(e.ExecuteQuery(`SELECT count(value) FROM cpu`, "db0", 0)); !reflect.DeepEqual(a, []*influxql.Result{
// 		{
// 			StatementID: 0,
// 			Err:         errors.New("max select series count exceeded: 4 series"),
// 		},
// 	}) {
// 		t.Fatalf("unexpected results: %s", spew.Sdump(a))
// 	}
// }

// // Ensure query executor can enforce a maximum bucket selection count.
// func TestQueryExecutor_ExecuteQuery_MaxSelectBucketsN(t *testing.T) {
// 	e := DefaultQueryExecutor()
// 	e.StatementExecutor.MaxSelectBucketsN = 3

// 	// The meta client should return a single shards on the local node.
// 	e.MetaClient.ShardsByTimeRangeFn = func(sources influxql.Sources, tmin, tmax time.Time) (a []meta.ShardInfo, err error) {
// 		return []meta.ShardInfo{
// 			{ID: 100, Owners: []meta.ShardOwner{{NodeID: 0}}},
// 		}, nil
// 	}

// 	var ic IteratorCreator
// 	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
// 		return &FloatIterator{
// 			Points: []influxql.FloatPoint{{Name: "cpu", Time: int64(0 * time.Second), Aux: []interface{}{float64(100)}}},
// 		}, nil
// 	}
// 	ic.FieldDimensionsFn = func(sources influxql.Sources) (fields, dimensions map[string]struct{}, err error) {
// 		return map[string]struct{}{"value": struct{}{}}, nil, nil
// 	}
// 	ic.SeriesKeysFn = func(opt influxql.IteratorOptions) (influxql.SeriesList, error) {
// 		return influxql.SeriesList{
// 			{Name: "cpu", Aux: []influxql.DataType{influxql.Float}},
// 		}, nil
// 	}
// 	e.TSDBStore.ShardIteratorCreatorFn = func(id uint64) influxql.IteratorCreator { return &ic }

// 	// Verify all results from the query.
// 	if a := ReadAllResults(e.ExecuteQuery(`SELECT count(value) FROM cpu WHERE time >= '2000-01-01T00:00:05Z' AND time < '2000-01-01T00:00:35Z' GROUP BY time(10s)`, "db0", 0)); !reflect.DeepEqual(a, []*influxql.Result{
// 		{
// 			StatementID: 0,
// 			Err:         errors.New("max select bucket count exceeded: 4 buckets"),
// 		},
// 	}) {
// 		t.Fatalf("unexpected results: %s", spew.Sdump(a))
// 	}
// }

// // QueryExecutor is a test wrapper for cluster.QueryExecutor.
// type QueryExecutor struct {
// 	*influxql.QueryExecutor

// 	MetaClient        MetaClient
// 	TSDBStore         TSDBStore
// 	StatementExecutor *cluster.StatementExecutor
// 	LogOutput         bytes.Buffer
// }

// // NewQueryExecutor returns a new instance of QueryExecutor.
// // This query executor always has a node id of 0.
// func NewQueryExecutor() *QueryExecutor {
// 	e := &QueryExecutor{
// 		QueryExecutor: influxql.NewQueryExecutor(),
// 	}
// 	e.StatementExecutor = &cluster.StatementExecutor{
// 		MetaClient: &e.MetaClient,
// 		TSDBStore:  &e.TSDBStore,
// 	}
// 	e.QueryExecutor.StatementExecutor = e.StatementExecutor

// 	var out io.Writer = &e.LogOutput
// 	if testing.Verbose() {
// 		out = io.MultiWriter(out, os.Stderr)
// 	}
// 	e.QueryExecutor.Logger = log.New(out, "[query] ", log.LstdFlags)

// 	return e
// }

// // DefaultQueryExecutor returns a QueryExecutor with a database (db0) and retention policy (rp0).
// func DefaultQueryExecutor() *QueryExecutor {
// 	e := NewQueryExecutor()
// 	e.MetaClient.DatabaseFn = DefaultMetaClientDatabaseFn
// 	return e
// }

// // ExecuteQuery parses query and executes against the database.
// func (e *QueryExecutor) ExecuteQuery(query, database string, chunkSize int) <-chan *influxql.Result {
// 	return e.QueryExecutor.ExecuteQuery(MustParseQuery(query), database, chunkSize, false, make(chan struct{}))
// }

var _ coordinator.TSDBStore = &TSDBStore{}

// TSDBStore is a mockable implementation of coordinator.TSDBStore.
type TSDBStore struct {
	CreateShardFn  func(database, policy string, shardID uint64, enabled bool) error
	WriteToShardFn func(shardID uint64, points []models.Point) error

	DeleteDatabaseFn        func(name string) error
	DeleteMeasurementFn     func(database, name string) error
	DeleteRetentionPolicyFn func(database, name string) error
	DeleteShardFn           func(id uint64) error
	DeleteSeriesFn          func(database string, sources []influxql.Source, condition influxql.Expr) error
	ShardIteratorCreatorFn  func(id uint64) influxql.IteratorCreator
	BackupShardFn           func(id uint64, since time.Time, w io.Writer) error
	MeasurementsFn          func(databse string, cond influxql.Expr) ([]string, error)
	RestoreShardFn          func(id uint64, r io.Reader) error
	TagValuesFn             func(database string, cond influxql.Expr) ([]tsdb.TagValues, error)
}

func (s *TSDBStore) CreateShard(database, policy string, shardID uint64, enabled bool) error {
	if s.CreateShardFn == nil {
		return nil
	}
	return s.CreateShardFn(database, policy, shardID, enabled)
}

func (s *TSDBStore) WriteToShard(shardID uint64, points []models.Point) error {
	return s.WriteToShardFn(shardID, points)
}

func (s *TSDBStore) RestoreShard(id uint64, r io.Reader) error {
	if s.RestoreShardFn == nil {
		return nil
	}

	return s.RestoreShard(id, r)
}

func (s *TSDBStore) BackupShard(id uint64, since time.Time, w io.Writer) error {
	if s.BackupShardFn == nil {
		return nil
	}
	return s.BackupShardFn(id, since, w)
}

func (s *TSDBStore) DeleteDatabase(name string) error {
	return s.DeleteDatabaseFn(name)
}

func (s *TSDBStore) DeleteMeasurement(database, name string) error {
	return s.DeleteMeasurementFn(database, name)
}

func (s *TSDBStore) DeleteRetentionPolicy(database, name string) error {
	return s.DeleteRetentionPolicyFn(database, name)
}

func (s *TSDBStore) DeleteSeries(database string, sources []influxql.Source, condition influxql.Expr) error {
	return s.DeleteSeriesFn(database, sources, condition)
}

func (s *TSDBStore) DeleteShard(id uint64) error {
	return s.DeleteShardFn(id)
}

func (s *TSDBStore) Measurements(database string, cond influxql.Expr) ([]string, error) {
	if s.MeasurementsFn == nil {
		return nil, nil
	}
	return s.MeasurementsFn(database, cond)
}

func (s *TSDBStore) TagValues(database string, cond influxql.Expr) ([]tsdb.TagValues, error) {
	if s.TagValuesFn == nil {
		return nil, nil
	}

	return s.TagValuesFn(database, cond)
}

// // MustParseQuery parses s into a query. Panic on error.
// func MustParseQuery(s string) *influxql.Query {
// 	q, err := influxql.ParseQuery(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return q
// }

// // ReadAllResults reads all results from c and returns as a slice.
// func ReadAllResults(c <-chan *influxql.Result) []*influxql.Result {
// 	var a []*influxql.Result
// 	for result := range c {
// 		a = append(a, result)
// 	}
// 	return a
// }

// // IteratorCreator is a mockable implementation of IteratorCreator.
// type IteratorCreator struct {
// 	CreateIteratorFn func(opt influxql.IteratorOptions) (influxql.Iterator, error)
// }

// type FieldMapper struct {
// 	FieldDimensionsFn func(sources influxql.Sources) (fields, dimensions map[string]struct{}, err error)
// }

// func (ic *IteratorCreator) CreateIterator(opt influxql.IteratorOptions) (influxql.Iterator, error) {
// 	return ic.CreateIteratorFn(opt)
// }

// func (ic *IteratorCreator) FieldDimensions(sources influxql.Sources) (fields, dimensions map[string]struct{}, err error) {
// 	return ic.FieldDimensionsFn(sources)
// }

// // FloatIterator is a represents an iterator that reads from a slice.
// type FloatIterator struct {
// 	Points []influxql.FloatPoint
// 	stats  influxql.IteratorStats
// }

// func (itr *FloatIterator) Stats() influxql.IteratorStats { return itr.stats }
// func (itr *FloatIterator) Close() error                  { return nil }

// // Next returns the next value and shifts it off the beginning of the points slice.
// func (itr *FloatIterator) Next() (*influxql.FloatPoint, error) {
// 	if len(itr.Points) == 0 {
// 		return nil, nil
// 	}

// 	v := &itr.Points[0]
// 	itr.Points = itr.Points[1:]
// 	return v, nil
// }
