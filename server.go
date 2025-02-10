package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"post-comment-system/graph"
	"post-comment-system/internal/repository"
	inmemory_repo "post-comment-system/internal/repository/inmemory"
	postgres2 "post-comment-system/internal/repository/postgres"
	"post-comment-system/internal/service/comment"
	"post-comment-system/internal/service/post"
	"post-comment-system/internal/service/subscriber_manager"
	inmemory_storage "post-comment-system/internal/storage/inmemory"
	"post-comment-system/internal/storage/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[main]: Не удалось загрузить .env файл, используем системные переменные")
	}
	storage := flag.String("storage", "inmemory", "Select storage: inmemory or postgres")
	flag.Parse()

	var postRepo repository.PostRepository
	var commentRepo repository.CommentRepository

	switch *storage {
	case "inmemory":
		str := inmemory_storage.NewInMemoryStorage()
		postRepo = inmemory_repo.NewInMemoryPostRepo(str)
		commentRepo = inmemory_repo.NewInMemoryCommentRepo(str)
		log.Println("connected to inmemory database")
		break
	case "postgres":
		db, err := postgres.NewDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		postRepo = postgres2.NewPostPostgresRepository(db)
		commentRepo = postgres2.NewPostgresCommentRepo(db)
		log.Println("connected to postgres database")
		break
	default:
		log.Fatalf("Unsupported storage type: %s", *storage)
	}

	sm := subscriber_manager.NewSubscriptionManager()

	postService := post.NewPostService(postRepo, commentRepo)
	commentService := comment.NewCommentService(commentRepo, sm)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		PostService:    postService,
		CommentService: commentService,
	}}))

	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
