# Cómo terminar la entrega

Copia los archivos de este ZIP sobre tu repositorio local conservando su carpeta `.git` y tu `.env`. Si no tienes `.env`, copia `.env.example` y ajusta los valores para Docker. No publiques `.env`.

1. Inicia Docker Desktop. Ejecuta `docker compose up -d --build` y `docker compose ps`.
2. Si ya existía el volumen MySQL, aplica 002 como indica `proyecto-final/verificacion.md` (no uses down -v).
3. Ejecuta las pruebas del documento de verificación y guarda las evidencias reales.
4. Realiza commits separados: análisis y ADR; migración y módulo; diagramas y pruebas; evidencias reales. Revisa `git status` antes de subir.
5. Ejecuta `git push` a tu fork y comparte el enlace con el docente.

Swagger original no fue regenerado; los endpoints nuevos están disponibles en requests.http y el script de pruebas. Si deseas actualizar Swagger, agrega anotaciones y ejecuta swag init.

Defensa: explicar el acoplamiento al libro (2 min), comparar alternativas y defender A (3 min), mostrar diagramas y endpoints (3 min), demostrar resultados reales (2 min). La decisión de siete días es una propuesta de negocio, no un requisito del docente.
