
https://nikiforovall.github.io/docker/2020/09/19/publish-package-to-ghcr.html

### montar el servicio en docker
```bash
docker compose up
```

### probar la api con el link de test
```bash
curl http://127.0.0.1:5002/
```

### probar la api con autenticacion basica para obtener el token
```bash
curl --location --request POST "http://127.0.0.1:5002/ObtenerToken" --header "Authorization: Basic dXN1YXJpbzpjb250cmFzZW5pYQ=="  
```

### probar la api con el token obtenido
```bash
curl --location --request POST "http://127.0.01:5002/ServicioConToken" \
--header "Authorization: Bearer 7ac84b55e6392bc512b65efac99e2be8.09b457cae987753781bb5c6c0c6de730.54ad170900899157f15167cad5985ad1" \
--header "Content-Type: application/json" \
--data-raw "{
  \"Nombre\": \"JORGE\",
  \"Documentos\": [
    {
      \"TipoDocumento\": \"DUI\",
      \"NumeroDocumento\": \"123456789-0\"
    },
    {
      \"TipoDocumento\": \"PASAPORTE\",
      \"NumeroDocumento\": \"A123456789\"
    }
  ] 
}"
```
