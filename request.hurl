# requests.hurl 
# TEST DE CADA ENDPOINT

# 1. createUsuario
POST http://localhost:8080/usuarios
Content-Type: application/json
{
    "Nombre": "Carloss",
    "Apellido": "Galardi",
    "Email": "carlos@gmail.com",
    "Contrasena": "fisica1"
}

HTTP 201
[Asserts]
jsonpath "$.IDUsuario" exists
jsonpath "$.Nombre" == "Carlos"
jsonpath "$.Apellido" == "Galardi"
jsonpath "$.Email" == "carlos@gmail.com"