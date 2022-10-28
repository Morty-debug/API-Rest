
### probar la api desde terminal

```bash

docker compose up

curl http://127.0.0.1:5002/

curl --location --request POST "http://127.0.0.1:5002/ObtenerToken" --header "Authorization: Basic dXN1YXJpbzpjb250cmFzZW5pYQ==" --header "Content-Type: application/json" 

```
