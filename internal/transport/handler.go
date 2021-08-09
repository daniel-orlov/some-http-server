package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"some-http-server/internal/types"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/myntra/golimit/store"
)

type QuoteSvc interface {
	Create(ctx context.Context, req *types.CreateQuoteRequestData) (*types.CreateQuoteResponseData, error)
	Read(ctx context.Context, req *types.GetQuoteRequestData) (*types.FullQuoteData, error)
}

type Handler struct {
	svc   QuoteSvc
	limit *store.Store
}

func NewHandler(svc QuoteSvc, limit *store.Store) *Handler {
	return &Handler{svc: svc, limit: limit}
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/quotes", h.create).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/quotes/{user:%s}/{id:%s}", "[0-9]+", "[0-9]+"), h.read).Methods(http.MethodGet)
}

type CreateQuoteRequest struct {
	Data types.CreateQuoteRequestData `json:"data"`
}

type CreateQuoteResponse struct {
	Data types.CreateQuoteResponseData `json:"data"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, "invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Validate request
	err := ValidateCreateQuoteRequest(&req)
	if err != nil {
		renderErrorResponse(w, "invalid request", http.StatusBadRequest)
		return
	}

	if blocked := h.limit.Incr(fmt.Sprint(req.Data.AccountID), 1, viper.GetInt32("limit_threshold"), viper.GetInt32("limit_window"), true); !blocked {
		renderErrorResponse(w, "limit exceeded", http.StatusTooManyRequests)
		return
	}

	res, err := h.svc.Create(r.Context(), &req.Data)
	if err != nil {
		renderErrorResponse(w, "create failed", http.StatusInternalServerError)
		return
	}

	renderResponse(w,
		&CreateQuoteResponse{Data: types.CreateQuoteResponseData{
			QuoteID:        res.QuoteID,
			TransactionFee: res.TransactionFee,
			EDT:            res.EDT,
		}},
		http.StatusCreated)
}

type GetQuoteResponse struct {
	Data types.FullQuoteData `json:"data"`
}

func (h *Handler) read(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		renderErrorResponse(w, "no id provided", http.StatusPreconditionFailed)
		return
	}

	quoteID, err := strconv.Atoi(id)
	if err != nil {
		renderErrorResponse(w, "invalid id", http.StatusPreconditionFailed)
		return
	}

	user, ok := mux.Vars(r)["user"]
	if !ok {
		renderErrorResponse(w, "no user id provided", http.StatusPreconditionFailed)
		return
	}

	accountID, err := strconv.Atoi(user)
	if err != nil {
		renderErrorResponse(w, "invalid user id", http.StatusPreconditionFailed)
		return
	}

	req := &types.GetQuoteRequestData{ID: uint64(quoteID), AccountID: uint64(accountID)}
	fqd, err := h.svc.Read(r.Context(), req)
	if err != nil {
		renderErrorResponse(w, "read failed", http.StatusNotFound)
		return
	}

	renderResponse(w,
		&GetQuoteResponse{
			Data: types.FullQuoteData{
				ID: fqd.ID,
				Req: &types.CreateQuoteRequestData{
					SourceCurrency: fqd.Req.SourceCurrency,
					TargetCurrency: fqd.Req.TargetCurrency,
					Amount:         fqd.Req.Amount,
					AccountID:      fqd.Req.AccountID,
				},
				Res: &types.CreateQuoteResponseData{
					TransactionFee: fqd.Res.TransactionFee,
					EDT:            fqd.Res.EDT,
				},
			},
		},
		http.StatusOK)
}
