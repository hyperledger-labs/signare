FROM python:3.12-slim

WORKDIR /app

ENV VENV_DIR=/var/cache/app/.venv

RUN apt-get update -y && apt-get install -y --no-install-recommends\
    git \
    python3-virtualenv

RUN virtualenv "$VENV_DIR"

RUN "$VENV_DIR/bin/pip" install --no-cache-dir \
    mkdocs-material \
    mkdocs-puml \
    mkdocs-render-swagger-plugin

COPY mkdocs/mkdocs.yml .
COPY open-api-specs/signare-api/index_bundled_openapi.yaml .

CMD ["/var/cache/app/.venv/bin/mkdocs", "serve", "--dev-addr=0.0.0.0:8000"]
