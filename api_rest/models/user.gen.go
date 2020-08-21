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

type User struct {
	ID           UserID         `json:"id" bunny:"id" `
	Username     string         `bunny:"username" json:"username" `
	Secret       string         `bunny:"secret" json:"secret" `
	ContactEmail string         `bunny:"contact_email" json:"contact_email" `
	Admin        bool           `bunny:"admin" json:"admin" `
	PhoneNumber  string         `bunny:"phone_number" json:"phone_number" `
	CreatedAt    time.Time      `json:"created_at" bunny:"created_at" `
	IsDeleted    bool           `bunny:"is_deleted" json:"is_deleted" `
	DeletedAt    _import00.Time `json:"deleted_at" bunny:"deleted_at" `
	R            *userR         `json:"-" toml:"-" yaml:"-"`
	L            userL          `json:"-" toml:"-" yaml:"-"`
}

var UserColumns = struct {
	ID           string
	Username     string
	Secret       string
	ContactEmail string
	Admin        string
	PhoneNumber  string
	CreatedAt    string
	IsDeleted    string
	DeletedAt    string
}{
	ID:           "id",
	Username:     "username",
	Secret:       "secret",
	ContactEmail: "contact_email",
	Admin:        "admin",
	PhoneNumber:  "phone_number",
	CreatedAt:    "created_at",
	IsDeleted:    "is_deleted",
	DeletedAt:    "deleted_at",
}

type userR struct {
	Rides           RideSlice
	CurrentVehicles VehicleSlice
	LastVehicles    VehicleSlice
}

type userL struct{}

var (
	userColumns              = []string{"id", "username", "secret", "contact_email", "admin", "phone_number", "created_at", "is_deleted", "deleted_at"}
	userPrimaryKeyColumns    = []string{"id"}
	userNonPrimaryKeyColumns = []string{"username", "secret", "contact_email", "admin", "phone_number", "created_at", "is_deleted", "deleted_at"}
)

type (
	UserSlice []*User

	userQuery struct {
		*queries.Query
	}
)

var (
	userType                 = reflect.TypeOf(&User{})
	userMapping              = queries.MakeStructMapping(userType)
	userPrimaryKeyMapping, _ = queries.BindMapping(userType, userMapping, userPrimaryKeyColumns)
	userInsertCacheMut       sync.RWMutex
	userInsertCache          = make(map[string]insertCache)
	userUpdateCacheMut       sync.RWMutex
	userUpdateCache          = make(map[string]updateCache)
)

func (q userQuery) One(ctx context.Context) (*User, error) {
	o := &User{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for user: %w", err)
	}

	return o, nil
}

func (q userQuery) First(ctx context.Context) (*User, error) {
	o := &User{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Errorf("models: failed to execute a one query for user: %w", err)
	}

	return o, nil
}

func (q userQuery) All(ctx context.Context) (UserSlice, error) {
	var o []*User

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Errorf("models: failed to assign all query results to User slice: %w", err)
	}

	return o, nil
}

func (q userQuery) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Errorf("models: failed to count user rows: %w", err)
	}

	return count, nil
}

func (q userQuery) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Errorf("models: failed to check if user exists: %w", err)
	}

	return count > 0, nil
}

func (o *User) Rides(mods ...qm.QueryMod) rideQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"user_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := Rides(queryMods...)
	queries.SetFrom(query.Query, "\"ride\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"ride\".*"})
	}

	return query
}

func (userL) LoadRides(ctx context.Context, slice []*User) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &userR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"user_id\" in (%s)",
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
			if local.ID == foreign.UserID {

				local.R.Rides = append(local.R.Rides, foreign)

			}
		}
	}

	return nil
}

func (o *User) CurrentVehicles(mods ...qm.QueryMod) vehicleQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"current_user_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := Vehicles(queryMods...)
	queries.SetFrom(query.Query, "\"vehicle\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"vehicle\".*"})
	}

	return query
}

func (userL) LoadCurrentVehicles(ctx context.Context, slice []*User) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &userR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"current_user_id\" in (%s)",
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
			if foreign.CurrentUserID.Valid && foreign.CurrentUserID.ID == local.ID {

				local.R.CurrentVehicles = append(local.R.CurrentVehicles, foreign)

			}
		}
	}

	return nil
}

func (o *User) LastVehicles(mods ...qm.QueryMod) vehicleQuery {
	queryMods := []qm.QueryMod{

		qm.Where("\"last_user_id\"=?", o.ID),
	}

	queryMods = append(queryMods, mods...)
	query := Vehicles(queryMods...)
	queries.SetFrom(query.Query, "\"vehicle\"")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"vehicle\".*"})
	}

	return query
}

func (userL) LoadLastVehicles(ctx context.Context, slice []*User) error {
	args := make([]interface{}, len(slice)*1)
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &userR{}
		}

		args[i*1+0] = obj.ID

	}

	where := fmt.Sprintf(
		"\"f\".\"last_user_id\" in (%s)",
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
			if foreign.LastUserID.Valid && foreign.LastUserID.ID == local.ID {

				local.R.LastVehicles = append(local.R.LastVehicles, foreign)

			}
		}
	}

	return nil
}

func Users(mods ...qm.QueryMod) userQuery {
	mods = append(mods, qm.From("\"user\""))
	return userQuery{NewQuery(mods...)}
}

func FindUser(ctx context.Context, id UserID, selectCols ...string) (*User, error) {
	userObj := &User{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM \"user\" WHERE \"id\"=$1", sel,
	)

	q := queries.Raw(query, id)

	err := q.Bind(ctx, userObj)
	if err != nil {
		return nil, errors.Errorf("models: unable to select from user: %w", err)
	}

	return userObj, nil
}

func (o *User) Insert(ctx context.Context, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no user provided for insertion")
	}

	var err error

	if len(whitelist) == 0 {
		whitelist = userColumns
	}

	key := makeCacheKey(whitelist)
	userInsertCacheMut.RLock()
	cache, cached := userInsertCache[key]
	userInsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping(userType, userMapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user\" (\"%s\") VALUES (%s)", strings.Join(whitelist, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user\" DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Errorf("models: unable to insert into user: %w", err)
	}

	if !cached {
		userInsertCacheMut.Lock()
		userInsertCache[key] = cache
		userInsertCacheMut.Unlock()
	}

	return nil
}

func (o *User) Update(ctx context.Context, whitelist ...string) error {
	var err error

	if len(whitelist) == 0 {
		whitelist = userNonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {

		return nil
	}

	key := makeCacheKey(whitelist)
	userUpdateCacheMut.RLock()
	cache, cached := userUpdateCache[key]
	userUpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE \"user\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, whitelist),
			strmangle.WhereClause("\"", "\"", len(whitelist)+1, userPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userType, userMapping, append(whitelist, userPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Errorf("models: unable to update user row: %w", err)
	}

	if !cached {
		userUpdateCacheMut.Lock()
		userUpdateCache[key] = cache
		userUpdateCacheMut.Unlock()
	}

	return nil
}

func (q userQuery) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to update all for user: %w", err)
	}

	return nil
}

func (o *User) Delete(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no User provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userPrimaryKeyMapping)
	sql := "DELETE FROM \"user\" WHERE \"id\"=$1"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete from user: %w", err)
	}

	return nil
}

func (q userQuery) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
		return errors.New("models: no userQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Errorf("models: unable to delete all from user: %w", err)
	}

	return nil
}

func (o UserSlice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Errorf("models: unable to delete all from user slice: %w", err)
	}

	return nil
}

func (o *User) Reload(ctx context.Context) error {
	ret, err := FindUser(ctx, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

func (o *UserSlice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	users := UserSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user\".* FROM \"user\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &users)
	if err != nil {
		return errors.Errorf("models: unable to reload all in UserSlice: %w", err)
	}

	*o = users

	return nil
}

func UserExists(ctx context.Context, id UserID) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user\" where \"id\"=$1 limit 1)"

	row := bunny.QueryRow(ctx, sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("models: unable to check if user exists: %w", err)
	}

	return exists, nil
}
