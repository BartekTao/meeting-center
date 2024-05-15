package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/BartekTao/nycu-meeting-room-api/internal/app"
	"github.com/BartekTao/nycu-meeting-room-api/internal/graph"
	"github.com/BartekTao/nycu-meeting-room-api/internal/graph/resolvers"
	infra "github.com/BartekTao/nycu-meeting-room-api/internal/infrastructure"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/auth"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/lock"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/middleware"
	"github.com/BartekTao/nycu-meeting-room-api/pkg/otelwrapper"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const defaultPort = "8080"

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := otelwrapper.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	mongoClient := infra.SetUpMongoDB()
	defer infra.ShutdownMongoDB(mongoClient)

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(mongoClient, client),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler(mongoClient *mongo.Client, rsClient *goredislib.Client) http.Handler {
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8888"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	jwtSecret := os.Getenv("JWT_KEY")
	if jwtSecret == "" {
		log.Fatal("You must set the JWT_KEY environment variable")
	}
	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	userRepo := infra.NewMongoUserRepo(mongoClient)
	authHandler := auth.NewGoogleOAuthHandler(userRepo)

	mux.HandleFunc("/auth/google/login", authHandler.Login)
	mux.HandleFunc("/auth/google/callback", authHandler.Callback)

	// Setup GraphQL server
	pool := goredis.NewPool(rsClient)
	rs := redsync.New(pool)
	locker := lock.NewRedsyncLocker(rs)
	roomRepo := infra.NewMongoRoomRepository(mongoClient)
	eventRepo := infra.NewMongoEventRepository(mongoClient)
	roomScheduleRepo := infra.NewRoomScheduleRepository(mongoClient)

	graphqlServer := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolvers.NewResolver(
			app.NewRoomService(roomRepo, roomScheduleRepo),
			app.NewEventService(eventRepo, locker),
			app.NewUserService(userRepo),
		),
	}))
	graphqlServer.AroundFields(tracer())
	graphqlServer.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, e)
		gqlErr.Message = "internal server error"
		return gqlErr
	})
	graphqlServer.Use(extension.FixedComplexityLimit(100))
	otelGraphqlHandler := otelhttp.NewHandler(graphqlServer, "GraphQL")
	jwtGraphqlHandler := jwtMiddleware(otelGraphqlHandler)

	// Wrap the GraphQL server with OpenTelemetry middleware
	otelHandler := otelhttp.WithRouteTag("/query", jwtGraphqlHandler)
	corsHandler := c.Handler(otelHandler)
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", corsHandler)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func tracer() graphql.FieldMiddleware {
	tracer := otel.Tracer("gqlgen-tracer")

	return func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		// Get the field name as operation name.
		rc := graphql.GetFieldContext(ctx)
		operationName := rc.Object + "." + rc.Field.Name

		// Start a new span.
		spanCtx, span := tracer.Start(ctx, operationName)
		defer span.End()

		// Continue execution to the next resolver.
		res, err = next(spanCtx)
		// Record any errors that occurred during resolution.
		if err != nil {
			span.RecordError(err)
		}

		return res, err
	}
}
