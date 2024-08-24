package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Foxtrot1388/MessaggioTask/internal/model"
	"github.com/Foxtrot1388/MessaggioTask/internal/service"
	"github.com/Foxtrot1388/MessaggioTask/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMessageCreate(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)

	wantresponse := []model.OutputMessage{
		{
			ID: 1,
		},
	}

	db := mocks.NewDbRepository(t)
	db.
		On("WithTr", mock.Anything, mock.Anything).
		Return(wantresponse, nil).
		Once()

	kafka := mocks.NewKafkaRepository(t)

	usercases := service.New(log, db, kafka)
	srvhttp := New(log, usercases)

	newmes := []string{"test"}
	example, _ := json.Marshal(newmes)
	req, err := http.NewRequest("POST", "/message/create", bytes.NewReader(example))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srvhttp.r.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusOK)

	body := rr.Body.String()
	var resp []model.OutputMessage
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, wantresponse, resp)

}
