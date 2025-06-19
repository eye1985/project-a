package user

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"project-a/internal/auth"
	"project-a/internal/consts"
	"project-a/internal/socket"
	"project-a/internal/testutil"
	"strings"
	"testing"
	"time"
)

func TestUserHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pgContainer, connStr := testutil.SetupTestContainer(ctx, t)

	t.Run(
		"Should update username", func(t *testing.T) {
			pool := testutil.CreateTestPoolAndCleanUp(t, ctx, connStr, pgContainer)
			ur := NewUserRepo(pool)
			ar := auth.NewRepo(pool)
			hub := socket.NewHub()
			go hub.Run()

			user, err := ur.InsertUser(ctx, "test", "test@test.com")
			require.NoError(t, err)

			sessionId := "test"
			_, err = ar.SetSession(
				ctx, &auth.SetSessionArgs{
					user.Id,
					sessionId,
					time.Now().Add(time.Hour * 1),
				},
			)
			require.NoError(t, err)
			newUserName := "newUserName"
			user.Username = newUserName
			body, err := json.Marshal(user)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			response := httptest.NewRecorder()
			sessionCtx := context.WithValue(request.Context(), consts.SessionCtxKey, []byte(sessionId))
			request = request.WithContext(sessionCtx)

			handler := NewUserHandler(ur, hub)
			handler.UpdateUserName(response, request)

			res := response.Result()
			defer res.Body.Close()
			require.Equal(t, res.StatusCode, http.StatusNoContent)

			user, err = ur.GetUserByEmail(ctx, user.Email)
			require.NoError(t, err)
			require.Equal(t, user.Username, newUserName)
		},
	)
}
