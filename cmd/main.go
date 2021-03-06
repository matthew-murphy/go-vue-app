package main

import (
  "log"
  "net/http"
  "os"
  web "github.com/matthew-murphy/kudo-oos/pkg/http"
  "github.com/matthew-murphy/kudo-oos/pkg/storage"
)

func main() {
  httpPort := os.Getenv("PORT")

  repo := storage.NewMongoRepository()
  webService := web.New(repo)

  log.Printf("Running on port %s\n", httpPort)
  log.Fatal(http.ListenAndServe(httpPort, webService.Router))
};
