
package main

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
	"bytes"
)

type Response struct {
	Estado        string
	Comparaciones []Matchs
}

type Matchs struct {
	Id            string
	Puntuacion    int
}

/*Ejemplo de JSON a USAR:
{
	"Estado":"Compatible",
	"Comparaciones":[
		{
			"Id":"Maria del Carmen",
			"Puntuacion":92
		},{
			"Id":"Maria Eva",
			"Puntuacion":80
		}
	]
}
*/

func main() {
	http.HandleFunc("/", writejson) //contruye y retorna JSON
	http.HandleFunc("/leer", readjson) //lee JSON desde una URL
	http.HandleFunc("/recibir", inputjson) //lee JSON cuando se lo envian via POST
	http.HandleFunc("/enviar", sendjson) //envia JSON via POST
	http.ListenAndServe(":8080", nil)
}

/* contruye y retorna JSON en localhost:8080 */
func writejson(w http.ResponseWriter, r *http.Request) {
	match0 := Matchs{"Maria del Carmen", 92}
	match1 := Matchs{"Maria Eva", 80}
	Matches := []Matchs{match0, match1}
	Estructura := Response {"Compatible",Matches}
	js, err := json.Marshal(Estructura)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"No se logro crear JSON\",\"Comparaciones\":null}") 
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

/* lee JSON desde la URL localhost:8080 para mostrarlo en localhost:8080/leer */
func readjson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	url := "http://localhost:8080/"
	res, err := http.Get(url)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"No se logro leer JSON de la URL "+url+"\",\"Comparaciones\":null}") 
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"No se logro leer JSON de la URL "+url+"\",\"Comparaciones\":null}") 
		return
	}
	var data Response
	json.Unmarshal(body, &data)
	fmt.Printf("Estado: %s\n", data.Estado)
	for i:=0; i<len(data.Comparaciones); i++{
		fmt.Printf("Id: %s Puntuacion: %d\n", data.Comparaciones[i].Id, data.Comparaciones[i].Puntuacion)
	}
	fmt.Fprintf(w, "JSON leido en consola\n") 
	return
}

/* lee JSON cuando se lo envian via POST a localhost:8080/recibir */
func inputjson(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Formato no compatible\",\"Comparaciones\":null}") 
		return
	}

	var data Response
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Estructura no compatible\",\"Comparaciones\":null}") 
		return
	}
	
	fmt.Printf("Estado: %s\n", data.Estado)
	for i:=0; i<len(data.Comparaciones); i++{
		fmt.Printf("Id: %s Puntuacion: %d\n", data.Comparaciones[i].Id, data.Comparaciones[i].Puntuacion)
	}
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"Estado\":\"JSON Recibido\",\"Comparaciones\":null}") 
	return
}

/* envia JSON via POST a localhost:8080/recibir y muestra resultados en localhost:8080/enviar */
func sendjson(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:8080/recibir"    
	
	match0 := Matchs{"Luisa", 100}
	match1 := Matchs{"Wendy", 75}
	Matches := []Matchs{match0, match1}
	Estructura := Response {"Compatible",Matches}
	js, err := json.Marshal(Estructura)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"Estado\":\"Imposible contruir JSON\",\"Comparaciones\":null}") 
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "Respuesta Estado: %s\n", resp.Status) 
	fmt.Fprintf(w, "Respuesta Encabezado: %s\n", resp.Header) 
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "Respuesta Cuerpo: %s\n", string(body)) 
	fmt.Fprintf(w, "JSON Enviado\n") 
	return
}
