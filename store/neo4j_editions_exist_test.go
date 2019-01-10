package store

import (
	"context"
	"fmt"
	"testing"

	dpbolt "github.com/ONSdigital/dp-bolt/bolt"
	"github.com/ONSdigital/dp-bolt/boltmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNeoDataStore_EditionExistsSuccess(t *testing.T) {
	Convey("given bolt.QueryForResult return count 1 for edition", t, func() {
		db := &boltmock.DB{
			QueryForResultFuncs: []boltmock.QueryFunc{
				func(query string, params map[string]interface{}, mapResult dpbolt.ResultMapper) error {
					return mapResult(
						&dpbolt.Result{
							Data: []interface{}{int64(1)},
						})
				},
			},
		}

		store := testStore
		store.bolt = db

		Convey("then EditionExists should return exists true and no error", func() {
			exists, err := store.EditionExists(context.Background(), testCodeListID, testEdition)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(db.QueryForResultCalls, ShouldHaveLength, 1)
			So(db.QueryForResultCalls[0].Query, ShouldEqual, fmt.Sprintf(countEditions, codeListLabel, testCodeListID, testEdition))
			So(db.QueryForResultCalls[0].Params, ShouldBeNil)
		})
	})
}

func TestNeoDataStore_EditionExistsMoreThanOneResult(t *testing.T) {
	Convey("given bolt.QueryForResult return count > 1 for edition", t, func() {
		db := &boltmock.DB{
			QueryForResultFuncs: []boltmock.QueryFunc{
				func(query string, params map[string]interface{}, mapResult dpbolt.ResultMapper) error {
					return mapResult(
						&dpbolt.Result{
							Data: []interface{}{int64(2)},
						})
				},
			},
		}

		store := testStore
		store.bolt = db

		Convey("then EditionExists should return exists true and no error", func() {
			exists, err := store.EditionExists(context.Background(), testCodeListID, testEdition)
			So(err.Error(), ShouldEqual, "editionExists: multiple editions found")
			So(exists, ShouldBeFalse)
			So(db.QueryForResultCalls, ShouldHaveLength, 1)
			So(db.QueryForResultCalls[0].Query, ShouldEqual, fmt.Sprintf(countEditions, codeListLabel, testCodeListID, testEdition))
			So(db.QueryForResultCalls[0].Params, ShouldBeNil)
		})
	})
}

func TestNeoDataStore_EditionExistsQueryForResultError(t *testing.T) {
	Convey("given bolt.QueryForResult returns an error", t, func() {
		db := &boltmock.DB{
			QueryForResultFuncs: []boltmock.QueryFunc{
				boltmock.ErrQueryFunc,
			},
		}

		store := testStore
		store.bolt = db

		Convey("then EditionExists should return exists false and the expected error", func() {
			exists, err := store.EditionExists(context.Background(), testCodeListID, testEdition)
			So(err, ShouldEqual, boltmock.Err)
			So(exists, ShouldBeFalse)
			So(db.QueryForResultCalls, ShouldHaveLength, 1)
			So(db.QueryForResultCalls[0].Query, ShouldEqual, fmt.Sprintf(countEditions, codeListLabel, testCodeListID, testEdition))
			So(db.QueryForResultCalls[0].Params, ShouldBeNil)
		})
	})
}