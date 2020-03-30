package main

import (
    "fmt"
    "time"
    "net/http"
    "cinema/home"
    "cinema/config"
    "github.com/gorilla/mux"
    "github.com/codegangsta/negroni"
)

// Route declare the name of the api, method used to get that api, pattern of url and function to use the specific api.
type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

// Routes is an array of Route types
type Routes []Route

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        home.IndexAction,
    },
    Route{
        "SaveBooking",
        "POST",
        "/save-booking",
        home.SaveBooking,
    },
    Route{
        "GetBooking",
        "POST",
        "/get-booking",
        home.GetBooking,
    },
    Route{
        "CancelBooking",
        "POST",
        "/cancel-booking",
        home.CancelBooking,
    },
}

func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w, r)

        msg := string(fmt.Sprintf(
            "%s %s %s %s %d %d %s %s IP: %s",
            r.RemoteAddr,
            r.Method,
            r.RequestURI,
            r.Proto,
            http.StatusOK,
            r.ContentLength,
            name,
            time.Since(start),
            r.Header.Get("x-forwarded-for"),
        ))
        fmt.Println(msg)
    })
}

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler

        handler = route.HandlerFunc
        handler = Logger(handler, route.Name)
        router.
        Methods(route.Method).
        Path(route.Pattern).
        Name(route.Name).
        Handler(handler)
    }
    return router
}

func main() {
    router := NewRouter()
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
    http.Handle("/", router)

    n := negroni.New()
    n.UseHandler(router)
    n.Run(":" + config.AppPort)
}
