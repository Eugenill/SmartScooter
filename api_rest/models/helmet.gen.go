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
	time "time"
)

type Helmet struct {
	ID           HelmetID     `bunny:"id" json:"id" `
	VehicleZone  VehicleZone  `bunny:"vehicle_zone" json:"vehicle_zone" `
	LastPing     time.Time    `bunny:"last_ping" json:"last_ping" `
	HelmetStatus HelmetStatus `bunny:"helmet_status" json:"helmet_status" `
	R            *helmetR     `json:"-" toml:"-" yaml:"-"`
	L            helmetL      `json:"-" toml:"-" yaml:"-"`
}

var HelmetColumns = struct {
	ID           string
	VehicleZone  string
	LastPing     string
	HelmetStatus string
}{
	ID:           "id",
	VehicleZone:  "vehicle_zone",
	LastPing:     "last_ping",
	HelmetStatus: "helmet_status",
}

type helmetR struct {
	Vehicles VehicleSlice
}

type helmetL struct{}

var (
	helmetColumns              = []string{"id", "vehicle_zone", "last_ping", "helmet_status"}
	helmetPrimaryKeyColumns    = []string{"id"}
	helmetNonPrimaryKeyColumns = []string{"vehicle_zone", "last_ping", "helmet_status"}
)

type (
	HelmetSlice []*Helmet

	helmetQuery struct {
		*queries.Query
	}
)

var (
	helmetType                 = reflect.TypeOf(&Helmet{})
	helmetMapping              = queries.MakeStructMapping(helmetType)
	helmetPrimaryKeyMapping, _ = queries.BindMapping(helmetType, helmetMapping, helmetPrimaryKeyColumns)
	helmetInsertCacheMut       sync.RWMutex
	helmetInsertCache          = make(map[string]insertCache)
	helmetUpdateCacheMut       sync.RWMutex
	helmetUpdateCache          = make(map[string]updateCache)
)

func (q helmetQuery) One(ctx context.Context) (*Helmet, error) {
	o := &Helmet{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for helmet: %w", err)
	}

	return o, nil
}

func (q helmetQuery) First(ctx context.Context) (*Helmet, error) {
	o := &Helmet{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for helmet: %w", err)
	}

	return o, nil
}

func (q helmetQuery) All(ctx context.Context) (HelmetSlice, error) {
	var o []*Helmet

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to Helmet slice: %w", err)
	}

	return o, nil
}

func (q helmetQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count helmet rows: %w", err)
	}

	return count, nil
}

func (q helmetQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if helmet exists: %w", err)
	}

	return count > 0, nil
}

func (o *Helmet) Vehicles(mods ...qm.QueryMod) vehicleQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"helmet_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := Vehicles(queryMods...)
	queries.SetFrom(query.Query, "\"vehicle\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"vehicle\".*"})
	}

	return query
}

func (helmetL) LoadVehicles(ctx context.Context, slice []*Helmet) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &helmetR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"helmet_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*1, 1, 1),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("\"vehicle\" AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*Vehicle
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice Vehicle: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ID == foreign.HelmetID {

				local.R.Vehicles = append(local.R.Vehicles, foreign)

			}
		}
	}

	return nil
}

func Helmets(mods ...qm.QueryMod) helmetQuery {
	mods = append(mods, qm.From("\"helmet\""))
	return helmetQuery{NewQuery(mods...)}
}

func FindHelmet(ctx context.Context, id HelmetID, selectCols ...string) (*Helmet, error) {
	helmetObj := &Helmet{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"helmet\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, helmetObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from helmet: %w", err)
	}

	return helmetObj, nil
}

func (o *Helmet) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no helmet provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = helmetColumns
	}

	key := makeCacheKey(whitelist)
	helmetInsertCacheMut.RLock()
	cache, cached := helmetInsertCache[key]
	helmetInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(helmetType, helmetMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"helmet\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"helmet\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into helmet: %w", err)
	}

	if !cached {
		helmetInsertCacheMut.Lock()
		helmetInsertCache[key] = cache
		helmetInsertCacheMut.Unlock()
	}

	return nil
}

func (o *Helmet) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = helmetNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	helmetUpdateCacheMut.RLock()
	cache, cached := helmetUpdateCache[key]
	helmetUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"helmet\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, helmetPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(helmetType, helmetMapping, append(whitelist, helmetPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update helmet row: %w", err)
	}

	if !cached {
		helmetUpdateCacheMut.Lock()
		helmetUpdateCache[key] = cache
		helmetUpdateCacheMut.Unlock()
	}

	return nil
}

func (q helmetQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for helmet: %w", err)
	}

	return nil
}

func (o *Helmet) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Helmet provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), helmetPrimaryKeyMapping)
	sql := "DELETE FROM \"helmet\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from helmet: %w", err)
	}

	return nil
}

func (q helmetQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no helmetQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from helmet: %w", err)
	}

	return nil
}

func (o HelmetSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Helmet slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), helmetPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"helmet\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, helmetPrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from helmet slice: %w", err)
	}

	return nil
}

func (o *Helmet) Reload(ctx context.Context) error {
	ret, err := FindHelmet(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *HelmetSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	helmets := HelmetSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), helmetPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"helmet\".* FROM \"helmet\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, helmetPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &helmets)
	if err != nil {
		return errors.Errorf("models: unable to reload all in HelmetSlice: %w", err)
	}

	*o = helmets

	return nil
}

func HelmetExists(ctx context.Context, id HelmetID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"helmet\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if helmet exists: %w", err)
	}

	return exists, nil
}
