package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const BASE_PATH = "/20230125"

func NewRouter(handler *EventServiceHandler) *mux.Router {
	var routes = []Route{
		{
			"Index",
			"GET",
			BASE_PATH + "/",
			Index,
		}, {
			"AddOrUpdateEvent",
			strings.ToUpper("Post"),
			BASE_PATH + "/events",
			handler.AddEvent,
		},
		{
			"AddOwner",
			strings.ToUpper("Put"),
			BASE_PATH + "/events/actions/share",
			handler.AddOwner,
		}, {
			"ListOWners",
			strings.ToUpper("GET"),
			BASE_PATH + "/eventsShared/{eventId}",
			handler.ListOwners,
		}, {
			"GetEvent",
			strings.ToUpper("Get"),
			BASE_PATH + "/events/{eventId}",
			handler.GetEvent,
		}, {
			"ListEvents",
			strings.ToUpper("Get"),
			BASE_PATH + "/events",
			handler.ListEvent,
		}, {
			"DeleteEvents",
			strings.ToUpper("Delete"),
			BASE_PATH + "/events/{eventId}",
			handler.DeleteEvent,
		},
		// Guests
		{
			"AddOrUpdateGuest",
			strings.ToUpper("Post"),
			BASE_PATH + "/guests",
			handler.AddGuest,
		}, {
			"GetGuest",
			strings.ToUpper("Get"),
			BASE_PATH + "/guests/{guestId}",
			handler.GetGuest,
		}, {
			"ListGuests",
			strings.ToUpper("Get"),
			BASE_PATH + "/guests",
			handler.ListGuest,
		}, {
			"DeleteGuest",
			strings.ToUpper("Delete"),
			BASE_PATH + "/guests/{guestId}",
			handler.DeleteGuest,
		},
		{
			"ActionCopyGuests",
			strings.ToUpper("Post"),
			BASE_PATH + "/actions/copy",
			handler.CopyGuests,
		},
		// Tasks
		{
			"AddOrUpdateTask",
			strings.ToUpper("Post"),
			BASE_PATH + "/tasks",
			handler.AddTask,
		}, {
			"GetTask",
			strings.ToUpper("Get"),
			BASE_PATH + "/tasks/{taskId}",
			handler.GetTask,
		}, {
			"ListTasks",
			strings.ToUpper("Get"),
			BASE_PATH + "/tasks",
			handler.ListTask,
		}, {
			"DeleteTask",
			strings.ToUpper("Delete"),
			BASE_PATH + "/tasks/{taskId}",
			handler.DeleteTask,
		},
		// Expenses
		{
			"AddOrUpdateExpense",
			strings.ToUpper("Post"),
			BASE_PATH + "/expenses",
			handler.AddExpense,
		}, {
			"GetExpense",
			strings.ToUpper("Get"),
			BASE_PATH + "/expenses/{expenseId}",
			handler.GetExpense,
		}, {
			"ListExpenses",
			strings.ToUpper("Get"),
			BASE_PATH + "/expenses",
			handler.ListExpenses,
		}, {
			"DeleteExpense",
			strings.ToUpper("Delete"),
			BASE_PATH + "/expenses/{expenseId}",
			handler.DeleteExpense,
		},
	}
	//
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = SetupGlobalMiddleware(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	// OPTIONS Method no op handler

	router.
		Methods("OPTIONS").
		Name("OptionsNoOp").
		PathPrefix(BASE_PATH).
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	return router
}

// Go Server API Index API!
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Events API is UP!")
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func SetupGlobalMiddleware(handler http.Handler, name string) http.Handler {
	return LoggerMiddleWare(JsonContentTypeMiddleWare(Authorization(Cors(handler))), name)
}

func Cors(inner http.Handler) http.Handler {
	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"*"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(inner)
}

// Authorization
func Authorization(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dangerous check  , but for now let's go with it.
		if IsLocal() {
			r.Header.Set("Authorization", "Bearer dummy")
		}
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			WriteError(w, 401, errors.New("authorization header not set "))
			return
		}

		inner.ServeHTTP(w, r)
	})
}

// Set application/json for all Responses in this Server
func JsonContentTypeMiddleWare(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
	})
}

// Log all Requests in this server
func LoggerMiddleWare(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func IsLocal() bool {
	lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT")
	return lambdaTaskRoot == ""
}
