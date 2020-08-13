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

type RideDetection struct {
	ID        RideDetectionID `bunny:"id" json:"id" `
	RideID    RideID          `bunny:"ride_id" json:"ride_id" `
	UserID    UserID          `bunny:"user_id" json:"user_id" `
	Detection Detection       `bunny:"detection__,bind" json:"detection" `
	R         *rideDetectionR `json:"-" toml:"-" yaml:"-"`
	L         rideDetectionL  `json:"-" toml:"-" yaml:"-"`
}

var RideDetectionColumns = struct {
	ID                         string
	RideID                     string
	UserID                     string
	DetectionTrafficLight      string
	DetectionObstacle          string
	DetectionTrafficSign       string
	DetectionLocationLatitude  string
	DetectionLocationLongitude string
	DetectionLocationAccuracy  string
	DetectionDetectedAt        string
}{
	ID:                         "id",
	RideID:                     "ride_id",
	UserID:                     "user_id",
	DetectionTrafficLight:      "detection__traffic_light",
	DetectionObstacle:          "detection__obstacle",
	DetectionTrafficSign:       "detection__traffic_sign",
	DetectionLocationLatitude:  "detection__location__latitude",
	DetectionLocationLongitude: "detection__location__longitude",
	DetectionLocationAccuracy:  "detection__location__accuracy",
	DetectionDetectedAt:        "detection__detected_at",
}

type rideDetectionR struct {
	Ride *Ride
	User *User
}

type rideDetectionL struct{}

var (
	rideDetectionColumns              = []string{"id", "ride_id", "user_id", "detection__traffic_light", "detection__obstacle", "detection__traffic_sign", "detection__location__latitude", "detection__location__longitude", "detection__location__accuracy", "detection__detected_at"}
	rideDetectionPrimaryKeyColumns    = []string{"id"}
	rideDetectionNonPrimaryKeyColumns = []string{"ride_id", "user_id", "detection__traffic_light", "detection__obstacle", "detection__traffic_sign", "detection__location__latitude", "detection__location__longitude", "detection__location__accuracy", "detection__detected_at"}
)

type (
	RideDetectionSlice []*RideDetection

	rideDetectionQuery struct {
		*queries.Query
	}
)

var (
	rideDetectionType                 = reflect.TypeOf(&RideDetection{})
	rideDetectionMapping              = queries.MakeStructMapping(rideDetectionType)
	rideDetectionPrimaryKeyMapping, _ = queries.BindMapping(rideDetectionType, rideDetectionMapping, rideDetectionPrimaryKeyColumns)
	rideDetectionInsertCacheMut       sync.RWMutex
	rideDetectionInsertCache          = make(map[string]insertCache)
	rideDetectionUpdateCacheMut       sync.RWMutex
	rideDetectionUpdateCache          = make(map[string]updateCache)
)

func (q rideDetectionQuery) One(ctx context.Context) (*RideDetection, error) {
	o := &RideDetection{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for ride_detection: %w", err)
	}

	return o, nil
}

func (q rideDetectionQuery) First(ctx context.Context) (*RideDetection, error) {
	o := &RideDetection{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for ride_detection: %w", err)
	}

	return o, nil
}

func (q rideDetectionQuery) All(ctx context.Context) (RideDetectionSlice, error) {
	var o []*RideDetection

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to RideDetection slice: %w", err)
	}

	return o, nil
}

func (q rideDetectionQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count ride_detection rows: %w", err)
	}

	return count, nil
}

func (q rideDetectionQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if ride_detection exists: %w", err)
	}

	return count > 0, nil
}

func (o *RideDetection) Ride(mods ...qm.QueryMod) rideQuery {
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

func (rideDetectionL) LoadRide(ctx context.Context, slice []*RideDetection) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &rideDetectionR{}
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

func (o *RideDetection) User(mods ...qm.QueryMod) userQuery {
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

func (rideDetectionL) LoadUser(ctx context.Context, slice []*RideDetection) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &rideDetectionR{}
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

func RideDetections(mods ...qm.QueryMod) rideDetectionQuery {
	mods = append(mods, qm.From("\"ride_detection\""))
	return rideDetectionQuery{NewQuery(mods...)}
}

func FindRideDetection(ctx context.Context, id RideDetectionID, selectCols ...string) (*RideDetection, error) {
	rideDetectionObj := &RideDetection{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"ride_detection\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, rideDetectionObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from ride_detection: %w", err)
	}

	return rideDetectionObj, nil
}

func (o *RideDetection) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no ride_detection provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = rideDetectionColumns
	}

	key := makeCacheKey(whitelist)
	rideDetectionInsertCacheMut.RLock()
	cache, cached := rideDetectionInsertCache[key]
	rideDetectionInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(rideDetectionType, rideDetectionMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"ride_detection\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"ride_detection\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into ride_detection: %w", err)
	}

	if !cached {
		rideDetectionInsertCacheMut.Lock()
		rideDetectionInsertCache[key] = cache
		rideDetectionInsertCacheMut.Unlock()
	}

	return nil
}

func (o *RideDetection) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = rideDetectionNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	rideDetectionUpdateCacheMut.RLock()
	cache, cached := rideDetectionUpdateCache[key]
	rideDetectionUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"ride_detection\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, rideDetectionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(rideDetectionType, rideDetectionMapping, append(whitelist, rideDetectionPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update ride_detection row: %w", err)
	}

	if !cached {
		rideDetectionUpdateCacheMut.Lock()
		rideDetectionUpdateCache[key] = cache
		rideDetectionUpdateCacheMut.Unlock()
	}

	return nil
}

func (q rideDetectionQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for ride_detection: %w", err)
	}

	return nil
}

func (o *RideDetection) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no RideDetection provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), rideDetectionPrimaryKeyMapping)
	sql := "DELETE FROM \"ride_detection\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from ride_detection: %w", err)
	}

	return nil
}

func (q rideDetectionQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no rideDetectionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from ride_detection: %w", err)
	}

	return nil
}

func (o RideDetectionSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no RideDetection slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), rideDetectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"ride_detection\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, rideDetectionPrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from rideDetection slice: %w", err)
	}

	return nil
}

func (o *RideDetection) Reload(ctx context.Context) error {
	ret, err := FindRideDetection(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *RideDetectionSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	rideDetections := RideDetectionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), rideDetectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"ride_detection\".* FROM \"ride_detection\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, rideDetectionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &rideDetections)
	if err != nil {
		return errors.Errorf("models: unable to reload all in RideDetectionSlice: %w", err)
	}

	*o = rideDetections

	return nil
}

func RideDetectionExists(ctx context.Context, id RideDetectionID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"ride_detection\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if ride_detection exists: %w", err)
	}

	return exists, nil
}
