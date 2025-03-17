package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

//go:embed web
var webFS embed.FS

type order struct {
	Name    string
	Drink   *drink
	Product *product
}
type drink struct {
	Name string
	Size int
}
type product struct {
	Name string

	Veggy bool
}

var orders = []*order{}

var mutex sync.Mutex

func main() {
	err := loadOrders()
	if err != nil {
		fmt.Println("Das laden der Bestellung ist fehlgeschlagen!", err)
		return
	}

	mux := http.NewServeMux()
	subFS, err := fs.Sub(webFS, "web")
	if err != nil {
		fmt.Println("Der Webpfad existiert nicht!", err)
		return
	}
	fileHandler := http.FileServer(http.FS(subFS))

	mux.Handle("GET /", fileHandler)
	mux.HandleFunc("GET /orders", getOrders)
	mux.HandleFunc("POST /orders", addOrder)
	mux.HandleFunc("DELETE /orders/{idx}", deleteorders)
	mux.HandleFunc("GET /orders/summary", summaryOrders)

	fmt.Println("Starte Server auf http://127.0.0.1:8080 ...")
	http.ListenAndServe("127.0.0.1:8080", mux)

}

func getOrders(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(resp)
	enc.Encode(orders)

}

func addOrder(resp http.ResponseWriter, req *http.Request) {

	dec := json.NewDecoder(req.Body)

	var newOrder order
	err := dec.Decode(&newOrder)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}

	if newOrder.Product == nil || newOrder.Product.Name == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Du musst ein produkt angeben!"))
		return
	}

	if newOrder.Name == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Du musst deinen Namen angeben!"))
		return
	}

	mutex.Lock()
	orders = append(orders, &newOrder)
	mutex.Unlock()

	err = saveOrders()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	resp.WriteHeader(http.StatusCreated)
}

func deleteorders(resp http.ResponseWriter, req *http.Request) {
	idxsrting := req.PathValue("idx")
	idx, err := strconv.Atoi(idxsrting)

	if err != nil {

		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return

	}
	fmt.Println("idx =", idx)
	if idx < 1 || idx > len(orders) {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("idx muss zwischen 1 und der Länge der Liste sein!"))
		return

	}

	mutex.Lock()
	orders = slices.Delete(orders, idx-1, idx)
	mutex.Unlock()

	err = saveOrders()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
}

type counts struct {
	Products map[string]int
	Drinks   map[string]int
}

func summaryOrders(resp http.ResponseWriter, req *http.Request) {
	countsProducts := map[string]int{}
	countsDrinks := map[string]int{}

	for _, order := range orders {
		key := strings.ToLower(order.Product.Name)
		if order.Product.Veggy {
			key += " veggy"
		} else {
			key += " notveggy"
		}

		countsProducts[key]++
		if order.Drink != nil {
			keyDrinks := strings.ToLower(order.Drink.Name)
			if order.Drink.Size == 2 {
				keyDrinks += " groß"
			} else if order.Drink.Size < 2 {
				keyDrinks += " klein"
			}
			countsDrinks[keyDrinks]++
		}

	}
	counts := counts{
		Products: countsProducts,
		Drinks:   countsDrinks,
	}

	resp.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(resp)
	enc.Encode(counts)

}

func saveOrders() error {
	file, err := os.Create("Orders.json")
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(orders)
	if err != nil {
		return err
	}
	return nil
}

func loadOrders() error {
	file, err := os.Open("Orders.json")
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&orders)
	if err != nil {
		return err
	}
	return nil
}
