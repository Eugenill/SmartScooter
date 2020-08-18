package models

import (
	context "context"
	fmt "fmt"
	errors "github.com/sqlbunny/errors"
	bunny "github.com/sqlbunny/sqlbunny/runtime/bunny"
	qm "github.com/sqlbunny/sqlbunny/runtime/qm"
	queries "github.com/sqlbunny/sqlbunny/runtime/queries"
	strmangle "github.com/sqlbunny/sqlbunny/runtime/strmangle"
	reflect "reflect"
	strings "strings"
	sync "sync"
)

type Path struct {
	ID     PathID `bunny:"id" json:"id" `
	RideID RideID `bunny:"ride_id" json:"ride_id" `
	Point  Point  `json:"point" bunny:"point__,bind" `
	R      *pathR `json:"-" toml:"-" yaml:"-"`
	L      pathL  `json:"-" toml:"-" yaml:"-"`
}

var PathColumns = struct {
	ID             string
	RideID         string
	PointLatitude  string
	PointLongitude string
	PointAccuracy  string
}{
	ID:             "id",
	RideID:         "ride_id",
	PointLatitude:  "point__latitude",
	PointLongitude: "point__longitude",
	PointAccuracy:  "point__accuracy",
}

type pathR struct {
	Ride  *Ride
	Rides RideSlice
}

type pathL struct{}

var (
	pathColumns              = []string{"id", "ride_id", "point__latitude", "point__longitude", "point__accuracy"}
	pathPrimaryKeyColumns    = []string{"id"}
	pathNonPrimaryKeyColumns = []string{"ride_id", "point__latitude", "point__longitude", "point__accuracy"}
)

type (
	PathSlice []*Path

	pathQuery struct {
		*queries.Query
	}
)

var (
	pathType                 = reflect.TypeOf(&Path{})
	pathMapping              = queries.MakeStructMapping(pathType)
	pathPrimaryKeyMapping, _ = queries.BindMapping(pathType, pathMapping, pathPrimaryKeyColumns)
	pathInsertCacheMut       sync.RWMutex
	pathInsertCache          = make(map[string]insertCache)
	pathUpdateCacheMut       sync.RWMutex
	pathUpdateCache          = make(map[string]updateCache)
)

func (q pathQuery) One(ctx context.Context) (*Path, error) {
	o := &Path{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for path: %w", err)
	}

	return o, nil
}

func (q pathQuery) First(ctx context.Context) (*Path, error) {
	o := &Path{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for path: %w", err)
	}

	return o, nil
}

func (q pathQuery) All(ctx context.Context) (PathSlice, error) {
	var o []*Path

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to Path slice: %w", err)
	}

	return o, nil
}

func (q pathQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count path rows: %w", err)
	}

	return count, nil
}

func (q pathQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if path exists: %w", err)
	}

	return count > 0, nil
}

func (o *Path) Ride(mods ...qm.QueryMod) rideQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"id\"=?", o.RideID),
	}

	queryMods = append(queryMods, mods...)
	query := Rides(queryMods...)
	queries.SetFrom(query.Query, "\"ride\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"ride\".*"})
	}

	return query
}

func (pathL) LoadRide(ctx context.Context, slice []*Path) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &pathR{}
		}

		args[i*1+0] = obj.RideID

	}

	where := fmt.Sprintf(
		"\"f\".\"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*1, 1, 1),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("\"ride\" AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*Ride
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice Ride: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.RideID == foreign.ID {

				local.R.Ride = foreign
				break

			}
		}
	}

	return nil
}

func (o *Path) Rides(mods ...qm.QueryMod) rideQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"path_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := Rides(queryMods...)
	queries.SetFrom(query.Query, "\"ride\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"ride\".*"})
	}

	return query
}

func (pathL) LoadRides(ctx context.Context, slice []*Path) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &pathR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"path_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*1, 1, 1),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("\"ride\" AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*Ride
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice Ride: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ID == foreign.PathID {

				local.R.Rides = append(local.R.Rides, foreign)

			}
		}
	}

	return nil
}

func Paths(mods ...qm.QueryMod) pathQuery {
	mods = append(mods, qm.From("\"path\""))
	return pathQuery{NewQuery(mods...)}
}

func FindPath(ctx context.Context, id PathID, selectCols ...string) (*Path, error) {
	pathObj := &Path{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"path\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, pathObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from path: %w", err)
	}

	return pathObj, nil
}

func (o *Path) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no path provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = pathColumns
	}

	key := makeCacheKey(whitelist)
	pathInsertCacheMut.RLock()
	cache, cached := pathInsertCache[key]
	pathInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(pathType, pathMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"path\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"path\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into path: %w", err)
	}

	if !cached {
		pathInsertCacheMut.Lock()
		pathInsertCache[key] = cache
		pathInsertCacheMut.Unlock()
	}

	return nil
}

func (o *Path) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = pathNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	pathUpdateCacheMut.RLock()
	cache, cached := pathUpdateCache[key]
	pathUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"path\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, pathPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(pathType, pathMapping, append(whitelist, pathPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update path row: %w", err)
	}

	if !cached {
		pathUpdateCacheMut.Lock()
		pathUpdateCache[key] = cache
		pathUpdateCacheMut.Unlock()
	}

	return nil
}

func (q pathQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for path: %w", err)
	}

	return nil
}

func (o *Path) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Path provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), pathPrimaryKeyMapping)
	sql := "DELETE FROM \"path\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from path: %w", err)
	}

	return nil
}

func (q pathQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no pathQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from path: %w", err)
	}

	return nil
}

func (o PathSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Path slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pathPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"path\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pathPrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from path slice: %w", err)
	}

	return nil
}

func (o *Path) Reload(ctx context.Context) error {
	ret, err := FindPath(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *PathSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	paths := PathSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pathPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"path\".* FROM \"path\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, pathPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &paths)
	if err != nil {
		return errors.Errorf("models: unable to reload all in PathSlice: %w", err)
	}

	*o = paths

	return nil
}

func PathExists(ctx context.Context, id PathID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"path\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if path exists: %w", err)
	}

	return exists, nil
}
