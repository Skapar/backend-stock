package database

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type colReplaceType int

const (
	RETURNING_OLD = "RETURNING *"
	RETURNING_NEW = "RETURNING %s"
	SELECT_OLD    = "SELECT *"
	SELECT_NEW    = "SELECT %s"
)

func (dbl *Database) GetWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error {
	cols := GetEntityDBFilds(returnValue)
	sql = ReplaceSelectWildcard(cols, sql)

	if _, ok := ctx.Value(CONTEXT_CACHE_MODE_KEY).(bool); ok {
		key, ok := ctx.Value(CONTEXT_CACHE_KEY).(string)
		if !ok {
			return fmt.Errorf(ERROR_KEY_EMPTY)
		}

		if dbl.cache.ExistKey(key) {
			return dbl.cache.Get(key, returnValue, false)
		} else {
			if err := pgxscan.Select(ctx, dbl.readPool, returnValue, sql, args...); err != nil { // get data from db
				return err
			}

			ttl := ctx.Value(CONTEXT_TTL_KEY).(time.Duration)
			if err := dbl.cache.Store(key, returnValue, ttl, false); err != nil { // store cache
				return err
			}

			return nil
		}
	}

	// without cache
	return pgxscan.Select(ctx, dbl.readPool, returnValue, sql, args...)
}

func (dbl *Database) GetOneWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error {
	cols := GetEntityDBFilds(returnValue)
	sql = ReplaceSelectWildcard(cols, sql)

	if _, ok := ctx.Value(CONTEXT_CACHE_MODE_KEY).(bool); ok {
		key, ok := ctx.Value(CONTEXT_CACHE_KEY).(string)
		if !ok {
			return fmt.Errorf(ERROR_KEY_EMPTY)
		}

		if dbl.cache.ExistKey(key) {
			return dbl.cache.Get(key, returnValue, false)
		} else {
			if err := pgxscan.Get(ctx, dbl.readPool, returnValue, sql, args...); err != nil { // get data from db
				if pgxscan.NotFound(err) {
					return fmt.Errorf("no rows")
				}
				return err
			}

			ttl := ctx.Value(CONTEXT_TTL_KEY).(time.Duration)
			if err := dbl.cache.Store(key, returnValue, ttl, false); err != nil { // store cache
				return err
			}

			return nil
		}
	}

	// without cache
	return pgxscan.Get(ctx, dbl.readPool, returnValue, sql, args...)
}

func (dbl *Database) InsertWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error {
	conn, err := dbl.writePool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if returnValue != nil {
		cols := GetEntityDBFilds(returnValue)
		sql = ReplaceReturningWildcard(cols, sql)

		dbl.log.Debug(sql)

		if err := pgxscan.Get(ctx, conn, returnValue, sql, args...); err != nil {
			return err
		}
	} else {
		cmd, err := conn.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		if !cmd.Insert() {
			dbl.log.Error(err)
		}
	}

	return nil
}

func (dbl *Database) UpdateWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error {
	conn, err := dbl.writePool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if returnValue != nil {
		cols := GetEntityDBFilds(returnValue)
		sql = ReplaceReturningWildcard(cols, sql)

		if err := pgxscan.Get(ctx, conn, returnValue, sql, args...); err != nil {
			return err
		}
	} else {
		cmd, err := conn.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		if !cmd.Update() {
			dbl.log.Error(err)
		}
	}

	return nil
}

func (dbl *Database) UpdateManyWithReplace(ctx context.Context, returnValue interface{}, sql string, args ...interface{}) error {
	conn, err := dbl.writePool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if returnValue != nil {
		cols := GetEntityDBFilds(returnValue)
		sql = ReplaceReturningWildcard(cols, sql)

		rows, err := conn.Query(ctx, sql, args...)
		if err != nil {
			return err
		}

		return pgxscan.ScanAll(returnValue, rows)
	} else {
		cmd, err := conn.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
		if !cmd.Update() {
			dbl.log.Error(err)
		}
	}

	return nil
}
