package main


import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


var conexion = "root:123456@tcp(127.0.0.1:3306)/ApiRest?charset=utf8"


type RespuestaToken struct {
	Token string
	Error int
}


func hasher(s string) []byte {
	val := sha256.Sum256([]byte(s))
	return val[:]
}



/*****************************************************/
/*
/* Autenticacion Basica
/*
/*****************************************************/
func authBasicHandler(handler http.HandlerFunc, userhash, passhash []byte, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare(hasher(user),userhash) != 1 || subtle.ConstantTimeCompare(hasher(pass),passhash) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			http.Error(w, "No autorizado.", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}



/*****************************************************/
/*
/* Autenticacion Token
/*
/*****************************************************/
func authTokenHandler(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var verificacion string
		authorization := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))
		db, err := sql.Open("mysql", conexion)
		defer db.Close()
		if err != nil {
			db.Close()
			fmt.Println("Error al aperturar la conexion a la base de datos")
			w.Header().Set("WWW-Authenticate", `Token realm="`+realm+`"`)
			http.Error(w, "No autorizado.", http.StatusUnauthorized)
			return
		}
		res, err := db.Query("SELECT ApiRest.validartocken(?)",idToken)
		defer res.Close()
		if err != nil {
			db.Close()
			res.Close()
			fmt.Println("Error al ejecutar la query")
			w.Header().Set("WWW-Authenticate", `Token realm="`+realm+`"`)
			http.Error(w, "No autorizado.", http.StatusUnauthorized)
			return
		}
		if res.Next() {
			err := res.Scan(&verificacion)
			if err != nil {
				db.Close()
				res.Close()
				fmt.Println("Error al obtener Token de la base de datos")
				w.Header().Set("WWW-Authenticate", `Token realm="`+realm+`"`)
				http.Error(w, "No autorizado.", http.StatusUnauthorized)
				return
			}
		} else {
			db.Close()
			res.Close()
			fmt.Println("Error al obtener Token de la base de datos")
			w.Header().Set("WWW-Authenticate", `Token realm="`+realm+`"`)
			http.Error(w, "No autorizado.", http.StatusUnauthorized)
			return
		}
		db.Close()
		res.Close()
		if verificacion == "Token Valido" {
			handler(w, r)
		} else {
			fmt.Println(verificacion)
			w.Header().Set("WWW-Authenticate", `Token realm="`+realm+`"`)
			http.Error(w, "No autorizado.", http.StatusUnauthorized)
			return
		}
	}
}



/*****************************************************/
/*
/* Metodo principal
/*
/*****************************************************/
func main() {
	userhash := hasher("usuario")
	passhash := hasher("contraseña")
	mux := http.NewServeMux()

	mux.HandleFunc("/GetToken", authBasicHandler(GetToken, userhash, passhash, "BasicAuth necesita credenciales"))
	mux.HandleFunc("/ServicioConToken", authTokenHandler(ServicioConToken, "TokenAuth necesita token"))

	server := &http.Server{ Addr: ":5002", Handler: mux, ReadTimeout: 60 * time.Second, WriteTimeout: 60 * time.Second }
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error de ListenAndServe")
	}
}



/*****************************************************/
/*
/* Servicio con Autenticacion Token
/*
/*****************************************************/
func ServicioConToken(w http.ResponseWriter, r *http.Request) {
	return
}



/*****************************************************/
/*
/* Servicio con Autenticacion Basica, Devuelve un Token
/*
/*****************************************************/
func GetToken(w http.ResponseWriter, r *http.Request) {
	var respuesta RespuestaToken
	db, err := sql.Open("mysql", conexion)
	defer db.Close()
	if err != nil {
		db.Close()
		respuesta.Error = 1
		js, _ := json.Marshal(respuesta)
		fmt.Println("Error al aperturar la conexion a la base de datos")
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	res, err := db.Query("SELECT ApiRest.creartocken()")
	defer res.Close()
	if err != nil {
		db.Close()
		res.Close()
		respuesta.Error = 1
		js, _ := json.Marshal(respuesta)
		fmt.Println("Error al ejecutar la query")
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	for res.Next() {
		err := res.Scan(&respuesta.Token)
		if err != nil {
			db.Close()
			res.Close()
			respuesta.Error = 1
			js, _ := json.Marshal(respuesta)
			fmt.Println("Error al obtener Token de la base de datos")
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
		respuesta.Error = 0
		fmt.Println("Token: (",respuesta.Token, ")")
	}
	db.Close()
	res.Close()
	js, _ := json.Marshal(respuesta)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}
