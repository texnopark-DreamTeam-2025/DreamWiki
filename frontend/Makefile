.PHONY: generate-client build generate-ssl-certs

generate-ssl-certs:
	@echo "Generating self-signed SSL certificates for development..."
	@mkdir -p ssl
	@openssl req -x509 -newkey rsa:4096 -keyout ssl/key.pem -out ssl/cert.pem -days 365 -nodes \
		-subj "/C=RU/ST=Moscow/L=Moscow/O=DreamWiki/OU=Development/CN=localhost" \
		-addext "subjectAltName=DNS:localhost,DNS:127.0.0.1,IP:127.0.0.1"
	@echo "SSL certificates generated in ssl/ directory"

generate-client:
	npm run openapi-ts

build: generate-client
	npm run build
