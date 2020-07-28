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

type Vehicle struct {
	ID          string      `bunny:"id" json:"id" `
	Token       string      `bunny:"token" json:"token" `
	NumberPlate string      `bunny:"number_plate" json:"number_plate" `
	VehicleZone VehicleZone `bunny:"vehicle_zone" json:"vehicle_zone" `
	HelmetID    HelmetID    `bunny:"helmet_id" json:"helmet_id" `
	R           *vehicleR   `json:"-" toml:"-" yaml:"-"`
	L           vehicleL    `json:"-" toml:"-" yaml:"-"`
}

var VehicleColumns = struct {
	ID          string
	Token       string
	NumberPlate string
	VehicleZone string
	HelmetID    string
}{
	ID:          "id",
	Token:       "token",
	NumberPlate: "number_plate",
	VehicleZone: "vehicle_zone",
	HelmetID:    "helmet_id",
}

type vehicleR struct {
	Helmet *Helmet
}

type vehicleL struct{}

var (
	vehicleColumns              = []string{"id", "token", "number_plate", "vehicle_zone", "helmet_id"}
	vehiclePrimaryKeyColumns    = []string{"id"}
	vehicleNonPrimaryKeyColumns = []string{"token", "number_plate", "vehicle_zone", "helmet_id"}
)

type (
	VehicleSlice []*Vehicle

	vehicleQuery struct {
		*queries.Query
	}
)

var (
	vehicleType                 = reflect.TypeOf(&Vehicle{})
	vehicleMapping              = queries.MakeStructMapping(vehicleType)
	vehiclePrimaryKeyMapping, _ = queries.BindMapping(vehicleType, vehicleMapping, vehiclePrimaryKeyColumns)
	vehicleInsertCacheMut       sync.RWMutex
	vehicleInsertCache          = make(map[string]insertCache)
	vehicleUpdateCacheMut       sync.RWMutex
	vehicleUpdateCache          = make(map[string]updateCache)
)

func (q vehicleQuery) One(ctx context.Context) (*Vehicle, error) {
	o := &Vehicle{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for vehicle: %w", err)
	}

	return o, nil
}

func (q vehicleQuery) First(ctx context.Context) (*Vehicle, error) {
	o := &Vehicle{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for vehicle: %w", err)
	}

	return o, nil
}

func (q vehicleQuery) All(ctx context.Context) (VehicleSlice, error) {
	var o []*Vehicle

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to Vehicle slice: %w", err)
	}

	return o, nil
}

func (q vehicleQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count vehicle rows: %w", err)
	}

	return count, nil
}

func (q vehicleQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if vehicle exists: %w", err)
	}

	return count > 0, nil
}

func (o *Vehicle) Helmet(mods ...qm.QueryMod) helmetQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"id\"=?", o.HelmetID),
	}

	queryMods = append(queryMods, mods...)
	query := Helmets(queryMods...)
	queries.SetFrom(query.Query, "\"helmet\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"helmet\".*"})
	}

	return query
}

func (vehicleL) LoadHelmet(ctx context.Context, slice []*Vehicle) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &vehicleR{}
		}

		args[i*1+0] = obj.HelmetID

	}

	where := fmt.Sprintf(
		"\"f\".\"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*1, 1, 1),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("\"helmet\" AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*Helmet
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice Helmet: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.HelmetID == foreign.ID {

				local.R.Helmet = foreign
				break

			}
		}
	}

	return nil
}

func Vehicles(mods ...qm.QueryMod) vehicleQuery {
	mods = append(mods, qm.From("\"vehicle\""))
	return vehicleQuery{NewQuery(mods...)}
}

func FindVehicle(ctx context.Context, id string, selectCols ...string) (*Vehicle, error) {
	vehicleObj := &Vehicle{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"vehicle\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, vehicleObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from vehicle: %w", err)
	}

	return vehicleObj, nil
}

func (o *Vehicle) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no vehicle provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = vehicleColumns
	}

	key := makeCacheKey(whitelist)
	vehicleInsertCacheMut.RLock()
	cache, cached := vehicleInsertCache[key]
	vehicleInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(vehicleType, vehicleMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"vehicle\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"vehicle\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into vehicle: %w", err)
	}

	if !cached {
		vehicleInsertCacheMut.Lock()
		vehicleInsertCache[key] = cache
		vehicleInsertCacheMut.Unlock()
	}

	return nil
}

func (o *Vehicle) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = vehicleNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	vehicleUpdateCacheMut.RLock()
	cache, cached := vehicleUpdateCache[key]
	vehicleUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"vehicle\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, vehiclePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(vehicleType, vehicleMapping, append(whitelist, vehiclePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update vehicle row: %w", err)
	}

	if !cached {
		vehicleUpdateCacheMut.Lock()
		vehicleUpdateCache[key] = cache
		vehicleUpdateCacheMut.Unlock()
	}

	return nil
}

func (q vehicleQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for vehicle: %w", err)
	}

	return nil
}

func (o *Vehicle) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Vehicle provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), vehiclePrimaryKeyMapping)
	sql := "DELETE FROM \"vehicle\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from vehicle: %w", err)
	}

	return nil
}

func (q vehicleQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no vehicleQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from vehicle: %w", err)
	}

	return nil
}

func (o VehicleSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no Vehicle slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehiclePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"vehicle\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehiclePrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from vehicle slice: %w", err)
	}

	return nil
}

func (o *Vehicle) Reload(ctx context.Context) error {
	ret, err := FindVehicle(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *VehicleSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	vehicles := VehicleSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehiclePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"vehicle\".* FROM \"vehicle\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehiclePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &vehicles)
	if err != nil {
		return errors.Errorf("models: unable to reload all in VehicleSlice: %w", err)
	}

	*o = vehicles

	return nil
}

func VehicleExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"vehicle\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if vehicle exists: %w", err)
	}

	return exists, nil
}
