rm *.pem
rm *.csr
rm *.key

# Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr -subj "/"

# Self signe pem certificate with CSR
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.pem
