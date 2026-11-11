# Trabajo Práctico - Programación Web
**Autores:** Carlos, Luca, Hermes

---

## Documentación de Persistencia y Dominio
La especificación detallada del modelo relacional, entidades, matriz de operaciones CRUD y reglas de negocio se encuentra en:
* [Documentación de Persistencia](docs/documentacion_dominio.md)

---

## Gestión y Ejecución del Servidor (Segundo Plano)

Para facilitar el desarrollo del frontend (HTMLs, CSS) y Server-Side Rendering (SSR), se dispone de un conjunto de comandos en el `Makefile` para controlar el servidor y su base de datos sin bloquear la terminal.

> **Nota previa:** Si se acaba de clonar el repositorio, asegúrate de otorgar permisos de ejecución a los scripts en la carpeta `tests/`:
> ```bash
> chmod +x tests/*.sh
> ```

### Comandos de Make

| Comando | Descripción |
| :--- | :--- |
| `make start` | **Levanta el entorno en segundo plano:**<br>1. Inicia el contenedor PostgreSQL (`docker compose up -d db`).<br>2. Espera a que la base de datos esté lista para aceptar conexiones.<br>3. Compila el binario de Go en `bin/server`.<br>4. Ejecuta el servidor en segundo plano (`http://localhost:8080`), redirigiendo la salida a `server.log` y guardando el PID en `bin/server.pid`. |
| `make stop` | **Cierre limpio del servidor:**<br>Envía una señal `SIGTERM` al proceso en segundo plano para un cierre exitoso y ordenado, liberando el puerto 8080 y limpiando el archivo PID.<br>*(Mantiene PostgreSQL corriendo para permitir reinicios rápidos).* |
| `make stop-all` | **Apagado total:**<br>Detiene el servidor Go y además apaga los contenedores Docker (`docker compose down`). |
| `make restart` | **Reinicio del servidor:**<br>Detiene el servidor actual y lo vuelve a compilar y arrancar (ideal tras modificar código Go o handlers). |
| `make status` | **Estado del sistema:**<br>Muestra si la base de datos PostgreSQL y el servidor Go están activos, indicando el puerto y el PID correspondiente. |
| `make logs` | **Visualización de logs en vivo:**<br>Muestra en tiempo real la consola del servidor (`tail -f server.log`) para depurar peticiones HTTP, renders de plantillas y errores. *(Salir con `Ctrl + C`)*. |

### Otros Comandos de Make

- `make test`: Ejecuta los tests unitarios, de integración y de API con Hurl (`./tests/test.sh`).
- `make clean`: Limpia y elimina los contenedores y volúmenes de Docker (`docker compose down -v --remove-orphans`).
- `make generate`: Regenera el código Go a partir de las consultas SQL utilizando `sqlc`.

---

### Ejecución Directa mediante Scripts en `tests/` (Opcional)

Si prefieres no usar `make`, puedes ejecutar directamente los scripts desde la carpeta `tests/`:

- Iniciar servidor: `./tests/start.sh` (o `./tests/server.sh start`)
- Detener servidor: `./tests/stop.sh` (o `./tests/server.sh stop`)
- Ver estado: `./tests/server.sh status`
- Ver logs: `./tests/server.sh logs`
