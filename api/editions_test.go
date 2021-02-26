package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	storetest "github.com/ONSdigital/dp-code-list-api/datastore/datastoretest"
	"github.com/ONSdigital/dp-code-list-api/models"

	dbmodels "github.com/ONSdigital/dp-graph/v2/models"

	"github.com/ONSdigital/dp-graph/v2/graph/driver"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// Edition models for testing
var (
	dbEdition1 = dbmodels.Edition{
		ID:    editionID1,
		Label: "label1",
	}

	dbEdition2 = dbmodels.Edition{
		ID:    editionID2,
		Label: "label2",
	}

	dbEditions = dbmodels.Editions{
		Items: []dbmodels.Edition{
			dbEdition1, dbEdition2,
		},
	}

	expectedEdition1 = models.Edition{
		ID:    editionID1,
		Label: "label1",
		Links: &models.EditionLinks{
			Self:     &models.Link{ID: editionID1, Href: fmt.Sprintf("%s/code-lists/%s/editions/%s", codeListURL, codeListID1, editionID1)},
			Editions: &models.Link{ID: "", Href: fmt.Sprintf("%s/code-lists/%s/editions", codeListURL, codeListID1)},
			Codes:    &models.Link{ID: "", Href: fmt.Sprintf("%s/code-lists/%s/editions/%s/codes", codeListURL, codeListID1, editionID1)},
		},
	}

	expectedEdition2 = models.Edition{
		ID:    editionID2,
		Label: "label2",
		Links: &models.EditionLinks{
			Self:     &models.Link{ID: editionID2, Href: fmt.Sprintf("%s/code-lists/%s/editions/%s", codeListURL, codeListID1, editionID2)},
			Editions: &models.Link{ID: "", Href: fmt.Sprintf("%s/code-lists/%s/editions", codeListURL, codeListID1)},
			Codes:    &models.Link{ID: "", Href: fmt.Sprintf("%s/code-lists/%s/editions/%s/codes", codeListURL, codeListID1, editionID2)},
		},
	}

	expectedEditions = models.Editions{
		Items:      []models.Edition{expectedEdition1, expectedEdition2},
		Count:      2,
		Offset:     0,
		Limit:      20,
		TotalCount: 2,
	}

	editionsPaginationTestOne = models.Editions{
		Items:      []models.Edition{expectedEdition1},
		Count:      1,
		Offset:     0,
		Limit:      1,
		TotalCount: 2,
	}

	editionsPaginationTestTwo = models.Editions{
		Items:      []models.Edition{expectedEdition2},
		Count:      1,
		Offset:     1,
		Limit:      7,
		TotalCount: 2,
	}

	editionsPaginationTestThree = models.Editions{
		Items:      []models.Edition{},
		Count:      0,
		Offset:     3,
		Limit:      1,
		TotalCount: 2,
	}

	editionsPaginationTestFour = models.Editions{
		Items:      []models.Edition{expectedEdition1},
		Count:      1,
		Offset:     0,
		Limit:      20,
		TotalCount: 1,
	}
)

func TestGetEditions(t *testing.T) {

	Convey("Get code list editions returns a status of http ok", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbEditions, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		validateBody(w.Body, &models.Editions{}, &expectedEditions)
	})

	Convey("Get code list editions returns a status of http not found if code list doesn't exist", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/12345/editions", codeListURL), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbmodels.Editions{}, driver.ErrNotFound
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("Get code list editions returns a status internal server error if store returns any other error", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/12345/editions", codeListURL), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbmodels.Editions{}, ErrInternal
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("Get code list editions returns a status internal server error if store returns an invalid edition (empty ID)", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/12345/editions", codeListURL), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbmodels.Editions{
					Items: []dbmodels.Edition{{ID: ""}},
				}, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("Get code list editions returns a status internal server error if writeBody fails", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbEditions, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.writeBody = failWriteBody
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

func TestGetEditions_Pagination(t *testing.T) {
	t.Parallel()

	Convey("When valid limit and offset query parameters are provided, then return editions information according to the offset and limit", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?offset=0&limit=1", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbEditions, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		validateBody(w.Body, &models.Editions{}, &editionsPaginationTestOne)
	})

	Convey("When valid limit above maximum and offset query parameters are provided, then return editions information according to the offset and limit", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?offset=1&limit=7", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbEditions, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		validateBody(w.Body, &models.Editions{}, &editionsPaginationTestTwo)
	})

	Convey("When offset value greater than count provided, then return zero items", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?offset=3&limit=1", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionsFunc: func(ctx context.Context, f string) (*dbmodels.Editions, error) {
				return &dbEditions, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		validateBody(w.Body, &models.Editions{}, &editionsPaginationTestThree)
	})

	Convey("When negative offset query parameters is provided, then 400 status returned", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?offset=-1", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		api := CreateCodeListAPI(mux.NewRouter(), &storetest.DataStoreMock{}, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When negative limit query parameters is provided, then 400 status returned", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?limit=-1", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		api := CreateCodeListAPI(mux.NewRouter(), &storetest.DataStoreMock{}, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When limit above default maximum is provided, then 400 status returned", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?limit=1001", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		api := CreateCodeListAPI(mux.NewRouter(), &storetest.DataStoreMock{}, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When non-integer query parameter value provided, then 400 status returned", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions?offset=x&limit=y", codeListURL, codeListID1), nil)
		w := httptest.NewRecorder()

		api := CreateCodeListAPI(mux.NewRouter(), &storetest.DataStoreMock{}, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestGetEdition(t *testing.T) {
	Convey("Get code list edition returns a status of http ok", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions/%s", codeListURL, codeListID1, editionID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionFunc: func(ctx context.Context, f, e string) (*dbmodels.Edition, error) {
				return &dbEdition1, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusOK)

		validateBody(w.Body, &models.Edition{}, &expectedEdition1)
	})

	Convey("Get code list edition returns a status of http not found if code list doesn't exist", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/12345/editions/2016", codeListURL), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionFunc: func(ctx context.Context, f, e string) (*dbmodels.Edition, error) {
				return &dbmodels.Edition{}, driver.ErrNotFound
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})

	Convey("Get code list edition returns a status internal server error if store returns any other error", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/12345/editions/2016", codeListURL), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionFunc: func(ctx context.Context, f, e string) (*dbmodels.Edition, error) {
				return &dbmodels.Edition{}, ErrInternal
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("Get code list edition returns a status internal server error if store returns any other error", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/12345/editions/2016", codeListURL), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionFunc: func(ctx context.Context, f, e string) (*dbmodels.Edition, error) {
				return &dbmodels.Edition{ID: ""}, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})

	Convey("Get code list edition returns a status internal server error if writeBody fails", t, func() {
		r := httptest.NewRequest("GET", fmt.Sprintf("%s/code-lists/%s/editions/%s", codeListURL, codeListID1, editionID1), nil)
		w := httptest.NewRecorder()

		mockDatastore := &storetest.DataStoreMock{
			GetEditionFunc: func(ctx context.Context, f, e string) (*dbmodels.Edition, error) {
				return &dbEdition1, nil
			},
		}

		api := CreateCodeListAPI(mux.NewRouter(), mockDatastore, codeListURL, datasetURL, defaultOffset, defaultLimit, maxLimit)
		api.writeBody = failWriteBody
		api.router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}
