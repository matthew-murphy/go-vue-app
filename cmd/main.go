package main

import (
  "log"
  "net/http"
  "os"
  web "github.com/{{ YOUR_GITHUB_USERNAME }}/kudo-oos/pkg/http"
  "github.com/{{ YOUR_GITHUB_USERNAME }}/kudo-oos/pkg/storage"
)

func main() {
  httpPort := os.Getenv("PORT")

  repo := storage.NewMongoRepository()
  webService := web.New(repo)

  log.Printf("Running on port %s\n", httpPort)
  log.Fatal(http.ListenAndServe(httpPort, webService.Router))
};
