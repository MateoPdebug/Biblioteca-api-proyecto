# Verificación del proyecto

**Estado:** pendiente de ejecución real en Docker. No se presentan resultados simulados.

## Preparación

Desde la raíz, ejecutar `docker compose up -d --build` y `docker compose ps`. Deben aparecer API, MySQL saludable y Adminer. Si se reutiliza el volumen de la práctica anterior, aplicar la migración sin eliminar datos (PowerShell):

```powershell
Get-Content migrations/002_laptops.sql -Raw | docker compose exec -T mysql mysql -uroot -proot_pass biblioteca
```

## Evidencia automática

```powershell
python scripts/verificar.py | Tee-Object -FilePath proyecto-final/evidencia-http.txt
```

El script registra usuarios de prueba con correos únicos, promueve uno a administrador en la base local, inicia sesión y comprueba los códigos HTTP y estados. No imprime tokens. Los registros de prueba quedan en la BD de desarrollo.

| Caso | Resultado esperado | Evidencia real |
|---|---|---|
| Docker build / ps | Tres servicios activos, MySQL healthy | Pendiente |
| health / ping | 200 y JSON esperado | Pendiente |
| Registro y login | 201 / 200 | Pendiente |
| Catálogo libros | 200 | Pendiente |
| Préstamo, historial y devolución libro | 201 / 200 / 200 | Pendiente |
| Crear laptop admin / estudiante | 201 / 403 | Pendiente |
| Catálogo laptop autenticado / sin token | 200 / 401 | Pendiente |
| Serial duplicado | 409 | Pendiente |
| Préstamo laptop | 201 | Pendiente |
| Segundo préstamo del mismo equipo | 409 | Pendiente |
| Historial propio, aislamiento | 200 | Pendiente |
| Devolución ajena / propia / repetida | 403 / 200 / 409 | Pendiente |
| Inventario disponible y nuevo préstamo | 200 / 201 | Pendiente |
| Devolución administrativa | 200 | Pendiente |

Adjuntar salida real del script y captura de `docker compose ps`; actualizar la tabla únicamente después de revisar los resultados. Si alguna prueba falla, conservar el error y corregirlo antes de entregar.

Prueba adicional: dos solicitudes simultáneas para el mismo serial deben producir un 201 y un 409. El script también devuelve el préstamo ganador.

Validación en el entorno de preparación: `gofmt` completado y sintaxis Python comprobada. `go test ./...` no pudo descargar las dependencias por falta de acceso DNS a proxy.golang.org. No hay Docker disponible en ese entorno; compilación y ejecución de la API quedan por verificar localmente.
