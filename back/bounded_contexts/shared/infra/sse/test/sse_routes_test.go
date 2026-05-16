package test

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	avataroptionusecases "brandtoonapi/bounded_contexts/creative_studio/avatar_option/useCases"
	authhttp "brandtoonapi/bounded_contexts/identity/auth/infra/http"
	sessiondomain "brandtoonapi/bounded_contexts/identity/session/domain"
	sessionmocks "brandtoonapi/bounded_contexts/identity/session/domain/mocks"
	userdomain "brandtoonapi/bounded_contexts/identity/user/domain"
	usermocks "brandtoonapi/bounded_contexts/identity/user/domain/mocks"
	sharedevents "brandtoonapi/bounded_contexts/shared/infra/events"
	sharedsse "brandtoonapi/bounded_contexts/shared/infra/sse"

	"github.com/go-chi/chi/v5"
)

func TestSSERouteRejectsUnauthenticatedRequests(t *testing.T) {
	t.Parallel()

	server := newSSEServer(t)
	request := httptest.NewRequest(http.MethodGet, "/events", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestSSERouteBroadcastsOnlyToMatchingUser(t *testing.T) {
	t.Parallel()

	eventBus := sharedevents.NewInMemoryEventBus()
	hub := sharedsse.NewHub()
	connector, err := sharedsse.NewConnector(
		eventBus,
		hub,
		[]string{avataroptionusecases.AvatarGenerationCompletedEventName},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	server := httptest.NewServer(newSSEServerWithConnector(t, connector))
	defer server.Close()

	client := server.Client()
	ctxUserOne, cancelUserOne := context.WithCancel(context.Background())
	defer cancelUserOne()
	requestUserOne, _ := http.NewRequestWithContext(ctxUserOne, http.MethodGet, server.URL+"/events", nil)
	requestUserOne.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v7"})
	responseUserOne, err := client.Do(requestUserOne)
	if err != nil {
		t.Fatalf("expected user one stream, got %v", err)
	}
	defer func() {
		_ = responseUserOne.Body.Close()
	}()

	ctxUserTwo, cancelUserTwo := context.WithCancel(context.Background())
	defer cancelUserTwo()
	requestUserTwo, _ := http.NewRequestWithContext(ctxUserTwo, http.MethodGet, server.URL+"/events", nil)
	requestUserTwo.AddCookie(&http.Cookie{Name: "brandtoon_session_id", Value: "session-v8"})
	responseUserTwo, err := client.Do(requestUserTwo)
	if err != nil {
		t.Fatalf("expected user two stream, got %v", err)
	}
	defer func() {
		_ = responseUserTwo.Body.Close()
	}()

	type readResult struct {
		event string
		err   error
	}

	readResultUserOne := make(chan readResult, 1)
	go func() {
		event, err := readSingleSSEEvent(responseUserOne.Body)
		readResultUserOne <- readResult{event: event, err: err}
	}()

	readResultUserTwo := make(chan readResult, 1)
	go func() {
		event, err := readSingleSSEEvent(responseUserTwo.Body)
		readResultUserTwo <- readResult{event: event, err: err}
	}()

	eventBus.Publish(avataroptionusecases.AvatarGenerationCompletedEvent{
		AvatarIDValue:   "avatar-v7",
		AvatarNameValue: "Studio Hero",
		OutcomeValue:    "SUCCESS",
		UserIDValue:     "user-v7",
	})

	select {
	case result := <-readResultUserOne:
		if result.err != nil {
			t.Fatalf("expected user one event, got %v", result.err)
		}
		if !strings.Contains(result.event, "avatar-generation.completed") {
			t.Fatalf("expected event name, got %q", result.event)
		}
		if !strings.Contains(result.event, "Studio Hero") {
			t.Fatalf("expected avatar name payload, got %q", result.event)
		}
		if !strings.Contains(result.event, `"outcome":"SUCCESS"`) {
			t.Fatalf("expected success outcome payload, got %q", result.event)
		}
		cancelUserOne()
	case <-time.After(time.Second):
		t.Fatalf("expected user one to receive event")
	}

	select {
	case result := <-readResultUserTwo:
		if result.err == nil {
			t.Fatalf("expected no event for user two, got %q", result.event)
		}
		if !errors.Is(result.err, context.Canceled) &&
			!strings.Contains(result.err.Error(), "use of closed network connection") {
			t.Fatalf("expected cancellation after no event, got %v", result.err)
		}
	case <-time.After(150 * time.Millisecond):
		cancelUserTwo()
	}
}

func newSSEServer(t *testing.T) http.Handler {
	t.Helper()

	eventBus := sharedevents.NewInMemoryEventBus()
	hub := sharedsse.NewHub()
	connector, err := sharedsse.NewConnector(
		eventBus,
		hub,
		[]string{avataroptionusecases.AvatarGenerationCompletedEventName},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	return newSSEServerWithConnector(t, connector)
}

func newSSEServerWithConnector(t *testing.T, connector *sharedsse.Connector) http.Handler {
	t.Helper()

	router := chi.NewMux()
	deps := authhttp.AuthMiddlewareDeps{
		SessionRepo: &sessionmocks.SessionRepositoryMock{
			FindActiveByIDFunc: func(ctx context.Context, id string) (*sessiondomain.Session, error) {
				switch id {
				case "session-v7":
					return &sessiondomain.Session{ID: id, UserID: "user-v7", ExpiresAt: time.Now().Add(time.Hour)}, nil
				case "session-v8":
					return &sessiondomain.Session{ID: id, UserID: "user-v8", ExpiresAt: time.Now().Add(time.Hour)}, nil
				default:
					return nil, nil
				}
			},
		},
		UserRepo: &usermocks.UserRepositoryMock{
			FindByIDFunc: func(ctx context.Context, id string) (*userdomain.User, error) {
				return &userdomain.User{ID: id, Email: id + "@example.com", Name: id}, nil
			},
		},
	}
	sharedsse.RegisterRoutes(
		router,
		sharedsse.RouteDependencies{Connector: connector},
		authhttp.AuthMiddleware(deps),
	)
	return router
}

func readSingleSSEEvent(body io.Reader) (string, error) {
	reader := bufio.NewReader(body)
	lines := make([]string, 0, 2)
	for len(lines) < 2 {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lines = append(lines, trimmed)
	}

	return strings.Join(lines, "\n"), nil
}
