# httpwrap
----------
[![Go Reference](https://pkg.go.dev/badge/github.com/apourchet/httpwrap.svg)](https://pkg.go.dev/github.com/apourchet/httpwrap)
![Unit Tests Workflow](https://github.com/apourchet/httpwrap/actions/workflows/unit-tests.yml/badge.svg)

`httpwrap` is a thin wrapper around the default http library that lets you compose handlers
and automatically inject outputs into the inputs of the next handler.

The idea is that you can write your http endpoints as you would regular functions, and have
`httpwrap` automatically populate those structures from the incoming http requests.

This goes against most of the popular Go libraries that rely on an opaque and type-unsafe `Context`
to pass the relevant information to the business logic of your application.

## Simple API Example
```go
type ListMoviesParams struct {
    // This will be populated from the cookies sent with the request.
    UserID string `http:"cookie=x-user-id"`
    
    // This will be populated from the URL query parameters that are 
    // sent with the request.
    ReleaseYear int `http:"query=release-year"`
    
    // The rest, by default, will be parsed from the body of the request 
    // interpreted as JSON.
    Director *string
    Actor string
}

type ListMoviesResponse struct {
    Movies []string `json:"movies"`
}

// The response and the error will automatically be written to the stdlib http.ResponseWriter. By
// default, the response will be JSON marshalled and the status code will be 200 OK.
func ListMovies(params ListMoviesParams) (ListMoviesResponse, error) {
    if params.ReleaseYear != 2022 {
        return httpwrap.NewHTTPError(http.StatusBadRequest, "Only 2022 movies are searchable.")
    }
	
    ...
		
    return ListMoviesResponse{
        Movies: []string{
            "Finding Nemo",
            "Good Will Hunting",
        },
    }
}
```

## Raw HTTP Access
For certain endpoints or applications, it can be desirable to forego the automatic sending of the response or 
error with JSON. The example below shows how this is done, which looks pretty much identical to vanilla Go:
```go
func RawHTTPHandler(rw http.ResponseWriter, req *http.Request) error { 
    if req.Origin != "..." {
        return httpwrap.NewHTTPError(http.StatusUnauthorized, "Bad origin.")
    }
    rw.WriteHeader(201)
    rw.Write([]byte{"Raw HTTP body here"})
    return nil
}
```

## Routing
Routing with `httpwrap` is nearly identical as you would otherwise do it. You can either use the standard lib
or any other routing libraries that you are used to. The following snippet uses `gorilla/mux`:
```go
import (
    "log"
    "net/http"

    "github.com/apourchet/httpwrap"
    "github.com/gorilla/mux"
)

func main() {
    // Tell the httpwrapper to run checkAPICreds as middleware before moving on to call
    // the endpoints themselves.
    httpWrapper := httpwrap.NewStandardWrapper().Before(checkAPICreds)

    // Using gorilla/mux for this example, but httpWrapper.Wrap will turn your regular endpoint
    // functions into the required http.HandlerFunc type.
    router := mux.NewRouter()
    router.Handle("/movies/list", httpWrapper.Wrap(ListMovies)).Methods("GET")
    router.Handle("/raw-handler", httpWrapper.Wrap(RawHTTPHandler)).Methods("GET")
    http.Handle("/", router)
	
    log.Fatal(http.ListenAndServe(":3000", router))
}

type APICredentials struct {
    Key string `http:"header=x-application-passcode"`
}

...

// checkAPICreds checks the api credentials passed into the request. Those APICredentials
// will be populated using the headers in the http request.
func checkAPICreds(creds APICredentials) error {
    if creds.Key == "my-secret-key" {
        return nil
    }
    return httpwrap.NewHTTPError(http.StatusForbidden, "Bad credentials.")
}

```
This example also displays a simple authorization middleware, `checkAPICreds`.

## Middleware
Middlewares can be used to either short-circuit the http request lifecycle and return early, or to provide additional 
information to the endpoint that gets called after it. The following example uses two separate middleware functions
to accomplish both.
```go
import (
    "log"
    "net/http"

    "github.com/apourchet/httpwrap"
    "github.com/gorilla/mux"
)

func main() {
    httpWrapper := httpwrap.NewStandardWrapper().
        Before(getUserAccountInfo)
    wrapperWithAccessCheck := httpWrapper.Before(ensureUserAccess)
	
    // The listMovies endpoint needs account information only, but the addMovies endpoint also performs
    // an access list check.
    router := mux.NewRouter()
    router.Handle("/movies/list", httpWrapper.Wrap(ListMovies)).Methods("GET")
    router.Handle("/movies/add", wrapperWithAccessCheck.Wrap(AddMovie)).Methods("PUT")
    http.Handle("/", router)

    log.Fatal(http.ListenAndServe(":3000", router))
}

type UserAuthenticationMaterial struct {
	BearerToken string `http:"header=Authorization"`
}

func getUserAccountInfo(authMaterial UserAuthenticationMaterial) (UserAccountInfo, error) {
    userId, err := decodeBearerToken(authMaterial.BearerToken)
    if err != nil {
        return httpwrap.NewHTTPError(http.StatusUnauthorized, "Bad authentication material.")
    }

    // Find the user information in the database for instance.
    accountInfo, err := database.FindUserInformation(userId)
    if err != nil {
        return err
    }
	
    ...
	
    return UserAccountInfo{
        UserID: userId,
        UserHasAccess: false,
    }, nil
}

func ensureUserAccess(accountInfo UserAccountInfo) error {
    if !accountInfo.UserHasAccess {
        // Returning an error from a middleware will short-circuit the rest of the request lifecycle and 
        // early-return this to the client.
        return httpwrap.NewHTTPError(http.StatusForbidden, "Access forbidden.")
    }
    return nil
}

// The two endpoints below will have access not only to the information provided by the middlewares, but
// can gather additional parameters from the http request by taking in extra arguments.
// NOTE: These two endpoints do not _have to_ take in the information from the middleware, so if accountInfo
// is not actually used in the endpoint, it can be omitted from the function signature altogether.
func ListMovies(accountInfo UserAccountInfo, params ListMoviesParams) (ListMoviesResponse, error) {
    ...
}

func AddMovie(accountInfo UserAccountInfo, params AddMovieParams) (AddMovieResponse, error) {
    ...
}
```