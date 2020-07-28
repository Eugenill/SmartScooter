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

type AllDetection struct {
	RideID    RideID         `bunny:"ride_id" json:"ride_id" `
	UserID    UserID         `json:"user_id" bunny:"user_id" `
	Detection Detection      `bunny:"detection__,bind" json:"detection" `
	R         *allDetectionR `json:"-" toml:"-" yaml:"-"`
	L         allDetectionL  `json:"-" toml:"-" yaml:"-"`
}

var AllDetectionColumns = struct {
	RideID                     string
	UserID                     string
	DetectionTrafficLight      string
	DetectionObstacle          string
	DetectionLocationLatitude  string
	DetectionLocationLongitude string
	DetectionLocationAccuracy  string
	DetectionDetectedAt        string
}{
	RideID:                     "ride_id",
	UserID:                     "user_id",
	DetectionTrafficLight:      "detection__traffic_light",
	DetectionObstacle:          "detection__obstacle",
	DetectionLocationLatitude:  "detection__location__latitude",
	DetectionLocationLongitude: "detection__location__longitude",
	DetectionLocationAccuracy:  "detection__location__accuracy",
	DetectionDetectedAt:        "detection__detected_at",
}

type allDetectionR struct {
	Ride *Ride
	User *User
}

type allDetectionL struct{}

var (
	allDetectionColumns              = []string{"ride_id", "user_id", "detection__traffic_light", "detection__obstacle", "detection__location__latitude", "detection__location__longitude", "detection__location__accuracy", "detection__detected_at"}
	allDetectionPrimaryKeyColumns    = []string{"ride_id"}
	allDetectionNonPrimaryKeyColumns = []string{"user_id", "detection__traffic_light", "detection__obstacle", "detection__location__latitude", "detection__location__longitude", "detection__location__accuracy", "detection__detected_at"}
)

type (
	AllDetectionSlice []*AllDetection

	allDetectionQuery struct {
		*queries.Query
	}
)

var (
	allDetectionType                 = reflect.TypeOf(&AllDetection{})
	allDetectionMapping              = queries.MakeStructMapping(allDetectionType)
	allDetectionPrimaryKeyMapping, _ = queries.BindMapping(allDetectionType, allDetectionMapping, allDetectionPrimaryKeyColumns)
	allDetectionInsertCacheMut       sync.RWMutex
	allDetectionInsertCache          = make(map[string]insertCache)
	allDetectionUpdateCacheMut       sync.RWMutex
	allDetectionUpdateCache          = make(map[string]updateCache)
)

func (q allDetectionQuery) One(ctx context.Context) (*AllDetection, error) {
	o := &AllDetection{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for all_detection: %w", err)
	}

	return o, nil
}

func (q allDetectionQuery) First(ctx context.Context) (*AllDetection, error) {
	o := &AllDetection{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for all_detection: %w", err)
	}

	return o, nil
}

func (q allDetectionQuery) All(ctx context.Context) (AllDetectionSlice, error) {
	var o []*AllDetection

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to AllDetection slice: %w", err)
	}

	return o, nil
}

func (q allDetectionQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count all_detection rows: %w", err)
	}

	return count, nil
}

func (q allDetectionQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if all_detection exists: %w", err)
	}

	return count > 0, nil
}

func (o *AllDetection) Ride(mods ...qm.QueryMod) rideQuery {
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

func (allDetectionL) LoadRide(ctx context.Context, slice []*AllDetection) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &allDetectionR{}
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

func (o *AllDetection) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"id\"=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)
	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"user\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"user\".*"})
	}

	return query
}

func (allDetectionL) LoadUser(ctx context.Context, slice []*AllDetection) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &allDetectionR{}
		}

		args[i*1+0] = obj.UserID

	}

	where := fmt.Sprintf(
		"\"f\".\"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*1, 1, 1),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("\"user\" AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*User
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice User: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {

				local.R.User = foreign
				break

			}
		}
	}

	return nil
}

func AllDetections(mods ...qm.QueryMod) allDetectionQuery {
	mods = append(mods, qm.From("\"all_detection\""))
	return allDetectionQuery{NewQuery(mods...)}
}

func FindAllDetection(ctx context.Context, rideID RideID, selectCols ...string) (*AllDetection, error) {
	allDetectionObj := &AllDetection{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"all_detection\" WHERE \"ride_id\"=$1", sel,
	)

	q := queries.Raw(query, rideID)

	err := q.Bind(ctx, allDetectionObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from all_detection: %w", err)
	}

	return allDetectionObj, nil
}

func (o *AllDetection) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no all_detection provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = allDetectionColumns
	}

	key := makeCacheKey(whitelist)
	allDetectionInsertCacheMut.RLock()
	cache, cached := allDetectionInsertCache[key]
	allDetectionInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(allDetectionType, allDetectionMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"all_detection\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"all_detection\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into all_detection: %w", err)
	}

	if !cached {
		allDetectionInsertCacheMut.Lock()
		allDetectionInsertCache[key] = cache
		allDetectionInsertCacheMut.Unlock()
	}

	return nil
}

func (o *AllDetection) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = allDetectionNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	allDetectionUpdateCacheMut.RLock()
	cache, cached := allDetectionUpdateCache[key]
	allDetectionUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"all_detection\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, allDetectionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(allDetectionType, allDetectionMapping, append(whitelist, allDetectionPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update all_detection row: %w", err)
	}

	if !cached {
		allDetectionUpdateCacheMut.Lock()
		allDetectionUpdateCache[key] = cache
		allDetectionUpdateCacheMut.Unlock()
	}

	return nil
}

func (q allDetectionQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for all_detection: %w", err)
	}

	return nil
}

func (o *AllDetection) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no AllDetection provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), allDetectionPrimaryKeyMapping)
	sql := "DELETE FROM \"all_detection\" WHERE \"ride_id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from all_detection: %w", err)
	}

	return nil
}

func (q allDetectionQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no allDetectionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from all_detection: %w", err)
	}

	return nil
}

func (o AllDetectionSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no AllDetection slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), allDetectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"all_detection\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, allDetectionPrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from allDetection slice: %w", err)
	}

	return nil
}

func (o *AllDetection) Reload(ctx context.Context) error {
	ret, err := FindAllDetection(ctx, o.RideID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *AllDetectionSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	allDetections := AllDetectionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), allDetectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"all_detection\".* FROM \"all_detection\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, allDetectionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &allDetections)
	if err != nil {
		return errors.Errorf("models: unable to reload all in AllDetectionSlice: %w", err)
	}

	*o = allDetections

	return nil
}

func AllDetectionExists(ctx context.Context, rideID RideID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"all_detection\" where \"ride_id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, rideID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if all_detection exists: %w", err)
	}

	return exists, nil
}
