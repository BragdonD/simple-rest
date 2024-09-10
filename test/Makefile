generate-ca:
	openssl req -x509 -new -nodes -keyout ./ca/ca.key -out ./ca/ca.crt -config ./ca/ca.conf -days 3650

generate-server-cert:
	openssl genrsa -out ./certs/server.key 2048
	openssl req -new -key ./certs/server.key -out ./certs/server.csr -config ./certs/server.conf
	openssl x509 -req -in ./certs/server.csr -CA ./ca/ca.crt -CAkey ./ca/ca.key -CAcreateserial -out ./certs/server.crt -days 365 -sha256 -extfile ./certs/server.conf -extensions v3_req

