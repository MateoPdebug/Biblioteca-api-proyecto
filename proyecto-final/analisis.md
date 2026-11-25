# Análisis del problema

La biblioteca actualmente presta libros. El nuevo requisito es incorporar laptops sin afectar el catálogo, la autenticación ni los préstamos existentes.

## Acoplamiento actual

- `migrations/001_init.sql`: `loans.book_id` es obligatorio y referencia `books`; las existencias se representan mediante `available_copies`.
- `internal/models/loan.go`: `Loan` y `LoanRequest` contienen `BookID`; `LoanDetail` contiene título y autor.
- `internal/repository/loan_repository.go`: las consultas de historial unen `loans` con `books` y la detección de duplicados utiliza `book_id`.
- `internal/service/loan_service.go`: `Borrow` consulta copias, rechaza duplicados, fija 14 días y decrementa inventario; `Return` incrementa copias.
- `internal/handler/loan_handler.go` y `routes.go`: las solicitudes de préstamo interpretan el recurso como libro. JWT y roles, en cambio, no dependen del recurso.

Agregar proyectores o tablets con el esquema actual exigiría repetir el modelo de inventario, persistencia y servicios o rediseñar el préstamo existente. Las salas además necesitarían reservas por intervalos.

## Opciones

| Opción | Cambios en BD | Duplicación y reutilización | Otro recurso | Riesgo |
|---|---|---|---|---|
| A. Dominio paralelo | Dos tablas nuevas, sin alterar libros | Duplica parte del flujo; reutiliza JWT, usuarios, conexión y patrón de capas | Nuevo módulo por tipo; costo medio-alto | Bajo para libros |
| B. Recurso generalizado | Tabla de recursos y préstamo común; migración de datos y referencias | Reutiliza flujo común; necesita políticas por tipo | Bajo para recursos similares; mayor para salas | Alto por migración y contratos REST |
| C. Dominios con servicio común | Tablas específicas y extracción de lógica compartida | Reutiliza reglas comunes mediante interfaces; conserva detalles de inventario | Medio; requiere adaptadores | Medio por refactorización de libros |

A ofrece alta cohesión interna y bajo acoplamiento entre dominios a cambio de repetición. B favorece extensibilidad cuando existen muchos tipos homogéneos, pero exige migrar el historial y preservar compatibilidad. C reduce duplicación, aunque introducir una abstracción con dos recursos puede resultar prematuro. El criterio prioritario para este alcance es evitar regresiones, no anticipar todos los recursos futuros.
