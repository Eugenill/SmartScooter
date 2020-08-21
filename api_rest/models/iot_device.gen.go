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
)

type IotDevice struct {
	ID              IotDeviceID    `bunny:"id" json:"id" `
	Name            string         `bunny:"name" json:"name" `
	LastPing        _import00.Time `bunny:"last_ping" json:"last_ping" `
	IotDeviceStatus DeviceStatus   `bunny:"iot_device_status" json:"iot_device_status" `
	R               *iotDeviceR    `json:"-" toml:"-" yaml:"-"`
	L               iotDeviceL     `json:"-" toml:"-" yaml:"-"`
}

var IotDeviceColumns = struct {
	ID              string
	Name            string
	LastPing        string
	IotDeviceStatus string
}{
	ID:              "id",
	Name:            "name",
	LastPing:        "last_ping",
	IotDeviceStatus: "iot_device_status",
}

type iotDeviceR struct {
	Vehicles VehicleSlice
}

type iotDeviceL struct{}

var (
	iotDeviceColumns              = []string{"id", "name", "last_ping", "iot_device_status"}
	iotDevicePrimaryKeyColumns    = []string{"id"}
	iotDeviceNonPrimaryKeyColumns = []string{"name", "last_ping", "iot_device_status"}
)

type (
	IotDeviceSlice []*IotDevice

	iotDeviceQuery struct {
		*queries.Query
	}
)

var (
	iotDeviceType                 = reflect.TypeOf(&IotDevice{})
	iotDeviceMapping              = queries.MakeStructMapping(iotDeviceType)
	iotDevicePrimaryKeyMapping, _ = queries.BindMapping(iotDeviceType, iotDeviceMapping, iotDevicePrimaryKeyColumns)
	iotDeviceInsertCacheMut       sync.RWMutex
	iotDeviceInsertCache          = make(map[string]insertCache)
	iotDeviceUpdateCacheMut       sync.RWMutex
	iotDeviceUpdateCache          = make(map[string]updateCache)
)

func (q iotDeviceQuery) One(ctx context.Context) (*IotDevice, error) {
	o := &IotDevice{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for iot_device: %w", err)
	}

	return o, nil
}

func (q iotDeviceQuery) First(ctx context.Context) (*IotDevice, error) {
	o := &IotDevice{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for iot_device: %w", err)
	}

	return o, nil
}

func (q iotDeviceQuery) All(ctx context.Context) (IotDeviceSlice, error) {
	var o []*IotDevice

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to IotDevice slice: %w", err)
	}

	return o, nil
}

func (q iotDeviceQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count iot_device rows: %w", err)
	}

	return count, nil
}

func (q iotDeviceQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if iot_device exists: %w", err)
	}

	return count > 0, nil
}

func (o *IotDevice) Vehicles(mods ...qm.QueryMod) vehicleQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"iot_device_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := Vehicles(queryMods...)
	queries.SetFrom(query.Query, "\"vehicle\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"vehicle\".*"})
	}

	return query
}

func (iotDeviceL) LoadVehicles(ctx context.Context, slice []*IotDevice) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &iotDeviceR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"iot_device_id\" in (%s)",
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
			if local.ID == foreign.IotDeviceID {

				local.R.Vehicles = append(local.R.Vehicles, foreign)

			}
		}
	}

	return nil
}

func IotDevices(mods ...qm.QueryMod) iotDeviceQuery {
	mods = append(mods, qm.From("\"iot_device\""))
	return iotDeviceQuery{NewQuery(mods...)}
}

func FindIotDevice(ctx context.Context, id IotDeviceID, selectCols ...string) (*IotDevice, error) {
	iotDeviceObj := &IotDevice{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"iot_device\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, iotDeviceObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from iot_device: %w", err)
	}

	return iotDeviceObj, nil
}

func (o *IotDevice) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no iot_device provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = iotDeviceColumns
	}

	key := makeCacheKey(whitelist)
	iotDeviceInsertCacheMut.RLock()
	cache, cached := iotDeviceInsertCache[key]
	iotDeviceInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(iotDeviceType, iotDeviceMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"iot_device\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"iot_device\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into iot_device: %w", err)
	}

	if !cached {
		iotDeviceInsertCacheMut.Lock()
		iotDeviceInsertCache[key] = cache
		iotDeviceInsertCacheMut.Unlock()
	}

	return nil
}

func (o *IotDevice) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = iotDeviceNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	iotDeviceUpdateCacheMut.RLock()
	cache, cached := iotDeviceUpdateCache[key]
	iotDeviceUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"iot_device\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, iotDevicePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(iotDeviceType, iotDeviceMapping, append(whitelist, iotDevicePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update iot_device row: %w", err)
	}

	if !cached {
		iotDeviceUpdateCacheMut.Lock()
		iotDeviceUpdateCache[key] = cache
		iotDeviceUpdateCacheMut.Unlock()
	}

	return nil
}

func (q iotDeviceQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for iot_device: %w", err)
	}

	return nil
}

func (o *IotDevice) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no IotDevice provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), iotDevicePrimaryKeyMapping)
	sql := "DELETE FROM \"iot_device\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from iot_device: %w", err)
	}

	return nil
}

func (q iotDeviceQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no iotDeviceQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from iot_device: %w", err)
	}

	return nil
}

func (o IotDeviceSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no IotDevice slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), iotDevicePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"iot_device\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, iotDevicePrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from iotDevice slice: %w", err)
	}

	return nil
}

func (o *IotDevice) Reload(ctx context.Context) error {
	ret, err := FindIotDevice(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *IotDeviceSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	iotDevices := IotDeviceSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), iotDevicePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"iot_device\".* FROM \"iot_device\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, iotDevicePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &iotDevices)
	if err != nil {
		return errors.Errorf("models: unable to reload all in IotDeviceSlice: %w", err)
	}

	*o = iotDevices

	return nil
}

func IotDeviceExists(ctx context.Context, id IotDeviceID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"iot_device\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if iot_device exists: %w", err)
	}

	return exists, nil
}
