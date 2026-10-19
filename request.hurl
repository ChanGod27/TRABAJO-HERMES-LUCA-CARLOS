# requests.hurl 
# TEST DE CADA ENDPOINT


#TABLA USUARIOS

# 1. createUsuario
POST http://localhost:8080/usuarios
Content-Type: application/json

{
    "Nombre": "Carlos",
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


# 2. getUserByEmail con el email carlos@gmail.com
GET http://localhost:8080/usuarios/email/carlos@gmail.com
HTTP 200
[Asserts]
jsonpath "$.IDUsuario" exists
jsonpath "$.Nombre" == "Carlos"
jsonpath "$.Apellido" == "Galardi"
jsonpath "$.Email" == "carlos@gmail.com"
jsonpath "$.Contrasena" == "fisica1"

# 3. getUserById con el id 1
GET http://localhost:8080/usuarios/1
HTTP 200
[Asserts]
jsonpath "$.IDUsuario" exists
jsonpath "$.Nombre" == "Carlos"
jsonpath "$.Apellido" == "Galardi"
jsonpath "$.Email" == "carlos@gmail.com"


#TABLA ListaTarea (MATERIAS)

# 1. createMateria 
POST http://localhost:8080/materias
Content-Type: application/json

{
    "IDUsuario": 1,
    "Titulo": "BaseDato2"
}

HTTP 201
[Asserts]
jsonpath "$.IDLista" exists
jsonpath "$.IDUsuario" == 1
jsonpath "$.Titulo" == "BaseDato2"


# 2. getListMateriasByUsuario obtenemos la materia del usuario 1
GET http://localhost:8080/materias/1
HTTP 200
[Asserts]
jsonpath "$" count >= 1
jsonpath "$[0].Titulo" == "BaseDato2"


# 3. putUpdateMateria
PUT http://localhost:8080/materias/1
Content-Type: application/json

{
    "IDLista": 1,
    "Titulo": "BaseDato2NombreCambiado"
}

HTTP 200
[Asserts]
jsonpath "$.IDLista" == 1
jsonpath "$.IDUsuario" == 1
jsonpath "$.Titulo" == "BaseDato2NombreCambiado"

#TABLA tarea 

# 1. createTarea
POST http://localhost:8080/tareas
Content-Type: application/json

{
    "IDLista": 1,
    "Titulo": "trabajoPracticoMongo",
    "FechaInicio": "2026-10-19T15:30:00Z",
    "FechaFinalizacion": "2026-11-19T15:30:00Z",
    "Descripcion": "Trabajo practico de prueba de bases NoSQL"
}

HTTP 201
[Asserts]
jsonpath "$.IDTarea" == 1 
jsonpath "$.IDLista" == 1 
jsonpath "$.Titulo" == "trabajoPracticoMongo"
jsonpath "$.FechaInicio" == "2026-10-19T15:30:00"
jsonpath "$.FechaFinalizacion" == "2026-11-19T15:30:00"
jsonpath "$.Descripcion" == "Trabajo practico de prueba de bases NoSQL"