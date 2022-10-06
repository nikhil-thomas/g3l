package lgwt_context

import (
    "context"
    "errors"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

type SpyStore struct {
    response string
    t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
    data := make(chan string, 1)
    go func() {
        var result string
        for _, c := range s.response {
            select {
            case <-ctx.Done():
                log.Println("spy store got cancelled")
                return
            default:
                time.Sleep(1 * time.Millisecond)
                result += string(c)
            }
        }
        data <- result
    }()
    select {
    case d := <-data:
        return d, nil
    case <-ctx.Done():
        log.Println("spy store got cancelled")
        return "", ctx.Err()
    }
}

type SpyResponseWriter struct {
    written bool
}

func (s *SpyResponseWriter) Header() http.Header {
    s.written = true
    return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
    s.written = true
    return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
    s.written = true
}

//func (s *SpyStore) Cancel() {
//    s.cancelled = true
//}

func TestServer(t *testing.T) {
    t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {

        data := "hello, world!"
        store := &SpyStore{response: data, t: t}
        svr := Server(store)

        request := httptest.NewRequest(http.MethodGet, "/", nil)
        cancellingCtx, cancel := context.WithCancel(request.Context())
        time.AfterFunc(5*time.Millisecond, cancel)
        request = request.WithContext(cancellingCtx)
        response := &SpyResponseWriter{}
        svr.ServeHTTP(response, request)
        if response.written {
            t.Error("a response should not have been written")
        }
    })

    t.Run("returns data from store", func(t *testing.T) {
        data := "hello, world"
        store := &SpyStore{response: data, t: t}
        svr := Server(store)
        request := httptest.NewRequest(http.MethodGet, "/", nil)
        response := httptest.NewRecorder()
        svr.ServeHTTP(response, request)
        //store.assertWasNotCancelled()

        if response.Body.String() != data {
            t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
        }

    })
}

//func (s *SpyStore) assertWasNotCancelled() {
//    s.t.Helper()
//    if s.cancelled {
//        s.t.Error("store was told to cancel")
//    }
//}
//
//func (s *SpyStore) assertWasCancelled() {
//    s.t.Helper()
//    if !s.cancelled {
//        s.t.Error("it should not have cancelled the store")
//    }
//}
