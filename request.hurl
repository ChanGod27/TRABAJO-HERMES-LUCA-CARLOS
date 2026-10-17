# requests.hurl 
# TEST DE CADA ENDPOINT

# 1. createUsuario
POST http://localhost:8080/usuarios
Content-Type: aplication/json
{
    "Nombre": "Carlos",
    "Apellido": "Galardi",
    "Email": "carlos@gmail.com",
    "Contrasena": "fisica1"
}
HTTP 201
[Asserts]
jsonpath "$.id_usuario" exists
jsonpath "$.Nombre" == "Carlos"
jsonpath "$.Apellido" == "Galardi"
jsonpath "$.Email" == "carlos@gmail.com"