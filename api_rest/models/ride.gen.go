package models

import (
	context "context"
	fmt "fmt"
	errors "github.com/sqlbunny/errors"
	bunny "github.com/sqlbunny/sqlbunny/runtime/bunny"
	qm "github.com/sqlbunny/sqlbunny/runtime/qm"
	queries "github.com/sqlbunny/sqlbunny/runtime/queries"
	strmangle "github.com/sqlbunny/sqlbunny/runtime/strmangle"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	reflect "reflect"
	strings "strings"
	sync "sync"
	time "time"
)

type Ride struct {
	ID         RideID         `bunny:"id" json:"id" `
	Distance   float32        `json:"distance" bunny:"distance" `
	StartedAt  time.Time      `json:"started_at" bunny:"started_at" `
	FinishedAt _import00.Time `bunny:"finished_at" json:"finished_at" `
	R          *rideR         `json:"-" toml:"-" yaml:"-"`
	L          rideL          `json:"-" toml:"-" yaml:"-"`
}

var RideColumns = struct {
	ID         string
	Distance   string
	StartedAt  string
	FinishedAt string
}{
	ID:         "id",
	Distance:   "distance",
	StartedAt:  "started_at",
	FinishedAt: "finished_at",
}

type rideR struct {
	AllDetection *AllDetection
}

type rideL struct{}

var (
	rideColumns              = []string{"id", "distance", "started_at", "finished_at"}
	ridePrimaryKeyColumns    = []string{"id"}
	rideNonPrimaryKeyColumns = []string{"distance", "started_at", "finished_at"}
)

type (
	RideSlice []*Ride

	rideQuery struct {
		*queries.Query
	}
)

var (
	rideType                 = reflect.TypeOf(&Ride{})
	rideMapping              = queries.MakeStructMapping(rideType)
	ridePrimaryKeyMapping, _ = queries.BindMapping(rideType, rideMapping, ridePrimaryKeyColumns)
	rideInsertCacheMut       sync.RWMutex
	rideInsertCache          = make(map[string]insertCache)
	rideUpdateCacheMut       sync.RWMutex
	rideUpdateCache          = make(map[string]updateCache)
)

func (q rideQuery) One(ctx context.Context) (*Ride, error) {
	o := &Ride{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for ride: %w", err)
	}

	return o, nil
}

func (q rideQuery) First(ctx context.Context) (*Ride, error) {
	o := &Ride{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for ride: %w", err)
	}

	return o, nil
}

func (q rideQuery) All(ctx context.Context) (RideSlice, error) {
	var o []*Ride

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to Ride slice: %w", err)
	}

	return o, nil
}

func (q rideQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count ride rows: %w", err)
	}

	return count, nil
}

func (q rideQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if ride exists: %w", err)
	}

	return count > 0, nil
}

func (o *Ride) AllDetection(mods ...qm.QueryMod) allDetectionQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"ride_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := AllDetections(queryMods...)
	queries.SetFrom(query.Query, "\"all_detection\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"all_detection\".*"})
	}

	return query
}

func (rideL) LoadAllDetection(ctx context.Context, slice []*Ride) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &rideR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"ride_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*1, 1, 1),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("\"all_detection\" AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*AllDetection
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice AllDetection: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ID == foreign.RideID {

				local.R.AllDetection = foreign
				break

			}
		}
	}

	return nil
}

func Rides(mods ...qm.QueryMod) rideQuery {
	mods = append(mods, qm.From("\"ride\""))
	return rideQuery{NewQuery(mods...)}
}

func FindRide(ctx context.Context, id RideID, selectCols ...string) (*Ride, error) {
	rideObj := &Ride{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"ride\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, rideObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from ride: %w", err)
	}

	return rideObj, nil
}

func (o *Ride) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no ride provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = rideColumns
	}

	key := makeCacheKey(whitelist)
	rideInsertCacheMut.RLock()
	cache, cached := rideInsertCache[key]
	rideInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(rideType, rideMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"ride\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"ride\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into ride: %w", err)
	}

	if !cached {
		rideInsertCacheMut.Lock()
		rideInsertCache[key] = cache
		rideInsertCacheMut.Unlock()
	}

	return nil
}

func (o *Ride) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = rideNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	rideUpdateCacheMut.RLock()
	cache, cached := rideUpdateCache[key]
	rideUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"ride\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, ridePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(rideType, rideMapping, append(whitelist, ridePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update ride row: %w", err)
	}

	if !cached {
		rideUpdateCacheMut.Lock()
		rideUpdateCache[key] = cache
		rideUpdateCacheMut.Unlock()
	}

	return nil
}

func (q rideQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for ride: %w", err)
	}

	return nil
}

func (o *Ride) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Ride provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), ridePrimaryKeyMapping)
	sql := "DELETE FROM \"ride\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from ride: %w", err)
	}

	return nil
}

func (q rideQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no rideQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from ride: %w", err)
	}

	return nil
}

func (o RideSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Ride slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ridePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"ride\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, ridePrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from ride slice: %w", err)
	}

	return nil
}

func (o *Ride) Reload(ctx context.Context) error {
	ret, err := FindRide(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *RideSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	rides := RideSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), ridePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"ride\".* FROM \"ride\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, ridePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &rides)
	if err != nil {
		return errors.Errorf("models: unable to reload all in RideSlice: %w", err)
	}

	*o = rides

	return nil
}

func RideExists(ctx context.Context, id RideID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"ride\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if ride exists: %w", err)
	}

	return exists, nil
}
