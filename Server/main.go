package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

var client *firestore.Client
var lastResponsePLL = time.Now()
var lastResponseOLL = time.Now()
var randPLL = 1
var randOLL = 22

func main() {
	var err error
	ctx := context.Background()
	client, err = firestore.NewClient(ctx, "rubiks-cube-api")
	if err != nil {
		log.Fatalf("Error initializing Cloud Firestore client: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/", rootHandler)
	r.HandleFunc("/v1/algorithm/{name}", algorithmHandler)
	r.HandleFunc("/v1/algorithmCategory/{category}", algorithmCategoryHandler)
	r.HandleFunc("/v1/randomPLL/", randomPLLHandler)
	r.HandleFunc("/v1/randomOLL/", randomOLLHandler)

	log.Println("Rubiks Algorithm REST API listening on port", port)
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Origin"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS", "PATCH", "CONNECT"}),
	)
	if err := http.ListenAndServe(":"+port, cors(r)); err != nil {
		log.Fatalf("Error launching Pets REST API server: %v", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{status: 'running'}")
}
func algorithmHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	ctx := context.Background()
	algorithm, err := getAlgorithm(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": '%s'}`, err)
		return
	}

	if algorithm.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		msg := fmt.Sprintf("`Algorithm \"%s\" not found`", name)
		fmt.Fprintf(w, fmt.Sprintf(`{"status": "fail", "data": {"title": %s}}`, msg))
		return
	}

	data, err := json.Marshal(algorithm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": "Unable to fetch algorithm: %s"}`, err)
		return
	}
	fmt.Fprintf(w, fmt.Sprintf(`%s`, data))
}

func randomPLLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	algorithm, err := getRandomPLL(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": '%s'}`, err)
		return
	}

	if algorithm.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		msg := fmt.Sprintf("`Error Fetching random Algorithm`")
		fmt.Fprintf(w, fmt.Sprintf(`{"status": "fail", "data": {"title": %s}}`, msg))
		return
	}

	data, err := json.Marshal(algorithm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": "Unable to fetch algorithm: %s"}`, err)
		return
	}
	fmt.Fprintf(w, fmt.Sprintf(`%s`, data))
}

func randomOLLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	algorithm, err := getRandomOLL(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": '%s'}`, err)
		return
	}

	if algorithm.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		msg := fmt.Sprintf("`Error Fetching random Algorithm`")
		fmt.Fprintf(w, fmt.Sprintf(`{"status": "fail", "data": {"title": %s}}`, msg))
		return
	}

	data, err := json.Marshal(algorithm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": "Unable to fetch algorithm: %s"}`, err)
		return
	}
	fmt.Fprintf(w, fmt.Sprintf(`%s`, data))
}

func algorithmCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	ctx := context.Background()
	algorithm, err := getAlgorithmByCategory(ctx, category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": '%s'}`, err)
		return
	}
	if len(algorithm) == 0 {
		w.WriteHeader(http.StatusNotFound)
		msg := fmt.Sprintf("`Algorithm Category \"%s\" not found`", category)
		fmt.Fprintf(w, fmt.Sprintf(`{"status": "fail", "data": {"title": %s}}`, msg))
		return
	}

	data, err := json.Marshal(algorithm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status": "fail", "data": "Unable to fetch algorithm: %s"}`, err)
		return
	}
	fmt.Fprintf(w, fmt.Sprintf(`%s`, data))
}

type Algorithm struct {
	Id         int    `firestore:"id"`
	Name       string `firestore:"name"`
	Moves      string `firestore:"moves"`
	VideoId    string `firestore:"videoid"`
	VideoStart int    `firestore:"videostart"`
	VideoEnd   int    `firestore:"videoend"`
	ShortNote  string `firestore:"shortnote"`
	Category   string `firestore:"category"`
	ImageUrl   string `firestore:"imageurl`
}

func getAlgorithm(ctx context.Context, name string) (*Algorithm, error) {
	query := client.Collection("Algorithms").Where("name", "==", name)
	iter := query.Documents(ctx)
	var c Algorithm
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		err = doc.DataTo(&c)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func getRandomPLL(ctx context.Context) (*Algorithm, error) {
	diffTime := time.Now().Sub(lastResponsePLL)
	if diffTime.Seconds() > 86400 {
		randPLL = rand.Intn(22) + 1
		lastResponsePLL = time.Now()
	}

	query := client.Collection("Algorithms").Where("Id", "==", randPLL)
	iter := query.Documents(ctx)
	var c Algorithm
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		err = doc.DataTo(&c)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func getRandomOLL(ctx context.Context) (*Algorithm, error) {
	diffTime := time.Now().Sub(lastResponseOLL)
	if diffTime.Seconds() > 86400 {
		randOLL = rand.Intn(57) + 23
		lastResponsePLL = time.Now()
	}

	query := client.Collection("Algorithms").Where("Id", "==", randPLL)
	iter := query.Documents(ctx)
	var c Algorithm
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		err = doc.DataTo(&c)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func getAlgorithmByCategory(ctx context.Context, category string) ([]Algorithm, error) {

	query := client.Collection("Algorithms").Where("category", "==", category)
	iter := query.Documents(ctx)

	algList := []Algorithm{}
	for {
		var c Algorithm
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		err = doc.DataTo(&c)
		algList = append(algList, c)
		if err != nil {
			return nil, err
		}
	}
	return algList, nil
}
