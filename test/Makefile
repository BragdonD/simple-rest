generate-ca:
	openssl req -x509 -new -nodes -keyout ./ca/ca.key -out ./ca/ca.crt -config ./ca/ca.conf -days 3650

generate-server-cert:
	openssl genrsa -out ./certs/server.key 2048
	openssl req -new -key ./certs/server.key -out ./certs/server.csr -config ./certs/server.conf
	openssl x509 -req -in ./certs/server.csr -CA ./ca/ca.crt -CAkey ./ca/ca.key -CAcreateserial -out ./certs/server.crt -days 365 -sha256 -extfile ./certs/server.conf -extensions v3_req

generate-client-cert:
	openssl genrsa -out ./certs/client.key 2048
	openssl req -new -key ./certs/client.key -out ./certs/client.csr -config ./certs/client.conf
	openssl x509 -req -in ./certs/client.csr -CA ./ca/ca.crt -CAkey ./ca/ca.key -CAcreateserial -out ./certs/client.crt -days 365 -sha256 -extfile ./certs/client.conf -extensions v3_req

generate-not-signed-client-crt:
	openssl genrsa -out ./certs/client-not-signed.key 2048
	openssl req -new -key ./certs/client-not-signed.key -out ./certs/client-not-signed.csr -config ./certs/client-not-signed.conf
	openssl x509 -req -in ./certs/client-not-signed.csr -key ./certs/client-not-signed.key -out ./certs/client-not-signed.crt -days 365 -sha256 -extfile ./certs/client-not-signed.conf -extensions v3_req


generate-all: generate-ca generate-server-cert generate-client-cert