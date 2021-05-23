package main

import (
	// "cerberus-security-laboratories/des-wristband-ui/internal/core"
	"cerberus-security-laboratories/des-wristband-ui/internal/gql"
	"cerberus-security-laboratories/des-wristband-ui/internal/gql/resolvers"
	"cerberus-security-laboratories/des-wristband-ui/server"

	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	// gqlHandler "github.com/99designs/gqlgen/handler"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/rs/cors"
)

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) { // Serves index.html
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}

var version = "1.0.0"

func main() {

	// Get Current Working Directory
	curr_wd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	// finalFolder := "Nurse"
	finalFolder := "Nurse"

	wd := filepath.Join(curr_wd, "internal", "core", "unittest_data", "TestWristbandFactory_NewWristband", finalFolder)

	// Handle command line arguments
	ipAddr := flag.String("ip", "localhost", "Server IP address. Deafult is localhost")
	port := flag.String("port", "8080", "Server port number. Default is 8080")
	dataPath := flag.String("data-path", wd, "Path to simulated wristband data files")
	dataPrefix := flag.String("data-prefix", "wbData_", "File name prefix for simulated wristband data files")
	tickPeriod := flag.Int("tick", 10000, "Tick period for simulated wristbands in milliseconds. Default is 20000")
	flag.Parse()

	r := mux.NewRouter()

	// Serving static files
	// through ui/build directory
	buildPath := path.Clean("ui/build")
	staticPath := path.Join(buildPath, "/static/")

	// GraphQL Executable Schema
	gqlExecutableSchema := gql.Config{Resolvers: resolvers.ResolverInitialisation(dataPath, dataPrefix, tickPeriod, finalFolder)}

	// websocket upgrader buffer size
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		// simple cors issue handler
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// CORS check for websocket connection
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	// GraphQL server

	// gqlSrv := gqlHandler.GraphQL(gql.NewExecutableSchema(gqlExecutableSchema), gqlHandler.WebsocketUpgrader(upgrader))
	gqlSrv := handler.New(gql.NewExecutableSchema(gqlExecutableSchema))
	// Add transport for websocket to GraphQL server to encode http request and decode graphql response
	gqlSrv.AddTransport(transport.Websocket{KeepAlivePingInterval: 10 * time.Second, Upgrader: upgrader})
	gqlSrv.AddTransport(transport.Options{})
	gqlSrv.AddTransport(transport.GET{})
	gqlSrv.AddTransport(transport.POST{})
	gqlSrv.AddTransport(transport.MultipartForm{})

	gqlSrv.SetQueryCache(lru.New(1000))

	gqlSrv.Use(extension.Introspection{})
	gqlSrv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// Serve static assets directly, via route /
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	r.Handle("/", server.Handler(buildPath))
	r.Use(c.Handler)
	// GraphQL API
	r.Handle("/gql-playground", playground.Handler("GraphQL playground", "/api"))
	// r.Handle("/gql-playground", gqlHandler.Playground("GraphQL playground", "/api"))
	r.Handle("/api", gqlSrv)

	// CORS for HTTP request polling
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
	)

	srv := &http.Server{
		Handler: handlers.CombinedLoggingHandler(os.Stdout, cors(r)),
		Addr:    *ipAddr + ":" + *port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Cerberus Wristband Server launched at http://%s:%s", *ipAddr, *port)
	log.Printf("\tConnect to http://%s:%s/gql-playground for GraphQL playground", *ipAddr, *port)
	log.Printf("\tWebsocket is triggered. Waiting for connection to the client...")
	log.Fatal(srv.ListenAndServe())
}
