package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/syucel96/simplebank/db/mock"
	db "github.com/syucel96/simplebank/db/sqlc"
	"github.com/syucel96/simplebank/token"
	"github.com/syucel96/simplebank/util"
)

func getRandomUser(t *testing.T) (user db.User) {
	user, _ = randomUser(t)
	return
}

func TestCreateTransferApi(t *testing.T) {
	users := []db.User{getRandomUser(t), getRandomUser(t), getRandomUser(t)}
	fromAccount := randomAccount(users[0].Username)
	toAccount := randomAccountWithCurrency(fromAccount.Currency)
	amount := util.RandomMoney(false, 1000)
	otherAccount := randomAccount(users[1].Username)
	for otherAccount.Currency == fromAccount.Currency {
		otherAccount = randomAccount(users[2].Username)
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)

				arg := db.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[2].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SameID",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   fromAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingFromAccountID",
			body: gin.H{
				"to_account_id": toAccount.ID,
				"amount":        amount,
				"currency":      fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[1].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingToAccountID",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingAmount",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingCurrency",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidFromAccountID",
			body: gin.H{
				"from_account_id": 0,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidToAccountID",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   "1",
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidAmount",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          "a",
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        "BTC",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          "-1.01",
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "AmountExceedsBalance",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          fmt.Sprintf("1%v", fromAccount.Balance),
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": otherAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[2].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(otherAccount.ID)).Times(1).Return(otherAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   otherAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(otherAccount.ID)).Times(1).Return(otherAccount, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountInternalServerError",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "ToAccountInternalServerError",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        fromAccount.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, users[0].Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)
				arg := db.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			urlPath := "/transfer"
			req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccountWithCurrency(currency string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(false),
		Currency: currency,
	}
}
