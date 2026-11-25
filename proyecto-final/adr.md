# ADR-001: Incorporar laptops como dominio paralelo

**Estado:** Aceptado

## Contexto

El préstamo existente depende directamente de libros y de sus copias. Una laptop necesita identificación individual, por lo que no conviene tratarla como una copia intercambiable. Consideramos el dominio paralelo (A), el recurso generalizado (B) y la extracción de servicios compartidos (C), detallados en el análisis.

## Decisión

Elegimos A. Se crean `laptops` y `laptop_loans` mediante la migración 002, con relaciones a usuarios e inventario. Las capas handler, service y repository conservan sus responsabilidades. Se integra un nuevo handler en main y se registran rutas separadas. No se altera la implementación de libros. Frente a B evitamos migrar datos y modificar contratos; frente a C aceptamos duplicar una parte del flujo para no refactorizar anticipadamente.

Cada laptop tiene serial único, marca, modelo y disponibilidad. Se presta durante siete días calendario, una política propuesta para facilitar la rotación de equipos escasos (no una regla institucional confirmada). Estudiantes y administradores autenticados pueden solicitar préstamos, como en libros; solo administradores crean inventario. El titular o un administrador puede devolver un préstamo. Esta última restricción mejora la autorización del nuevo módulo sin modificar el anterior.

El servicio valida disponibilidad y estado; usa una transacción para guardar el préstamo y cambiar inventario, o revertir ambos cambios. Un bloqueo de fila y un índice único sobre el identificador activo evitan doble asignación incluso ante concurrencia. La devolución cierra el préstamo y libera el equipo en otra transacción. El vencimiento no libera automáticamente el equipo: debe devolverse físicamente.

## Consecuencias

- Positivas: menor riesgo de regresión, cohesión por recurso, trazabilidad por serial, reutilización de JWT y conexión, consistencia transaccional.
- Negativas: más tablas y código; políticas duplicadas; los reportes globales requerirán consultar ambos dominios. Si se incorporan numerosos recursos, reevaluaremos B o C mediante un nuevo ADR.
- Operación: instalaciones con volumen MySQL existente deben ejecutar 002 manualmente; las nuevas lo hacen automáticamente. No se elimina información anterior.
