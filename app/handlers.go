package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nickduskey/api-demo/products"
	"github.com/nickduskey/api-demo/users"
	"github.com/nickduskey/api-demo/utils"
)

// NotImplementedHandler a placeholder
func (a *App) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, nil)
}

// StatusHandler provides a status endpoint
func (a *App) StatusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{"message": "API up and running!"}
	utils.RespondWithJSON(w, http.StatusOK, status)
}

// LoginHandler handles login auth
func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch user based on username
	// Hash request password from the body
	// Compare password
	// If valid user and good password generate the token
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	// create claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(signingKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	utils.RespondWithJSON(w, http.StatusOK, signedString)
}

// GetProduct responds with JSON product
func (a *App) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u64id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	id := uint(u64id)

	p := products.Product{ID: id}
	if err := p.GetSingle(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, p)
}

// GetProducts responds with JSON products
func (a *App) GetProducts(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	offset, _ := strconv.Atoi(v.Get("offset"))
	limit, _ := strconv.Atoi(v.Get("limit"))

	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, err := products.Get(a.DB, offset, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

// CreateProduct responds with JSON product created
func (a *App) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p products.Product
	decoder := json.NewDecoder(r.Body)
	fmt.Println(r.Body)
	if err := decoder.Decode(&p); err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.Create(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, p)
}

// UpdateProduct responds with updated product JSON
func (a *App) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u64id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	id := uint(u64id)

	var p products.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.Update(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, p)
}

// DeleteProduct responds with delete result JSON
func (a *App) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u64id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	id := uint(u64id)

	p := products.Product{ID: id}
	if err := p.Delete(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// GetUser responds with user JSON
func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u64id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	id := uint(u64id)

	u := users.User{ID: id}
	if err := u.GetSingle(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, u)
}

// GetUsers responds with users JSON
func (a *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	offset, _ := strconv.Atoi(v.Get("offset"))
	limit, _ := strconv.Atoi(v.Get("limit"))

	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := users.Get(a.DB, offset, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, users)
}

// CreateUser responds with created user JSON
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u users.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := u.Create(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, u)
}

// UpdateUser responds with updated user JSON
func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u64id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	id := uint(u64id)

	var u users.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.Update(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, u)
}

// DeleteUser responds with delete status JSON
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u64id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	id := uint(u64id)

	u := users.User{ID: id}
	if err := u.Delete(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
