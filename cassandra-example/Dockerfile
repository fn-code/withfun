FROM cassandra:3.11.6

RUN sed -i 's/AllowAllAuthenticator/PasswordAuthenticator/' /etc/cassandra/cassandra.yaml
# RUN cqlsh -u cassandra -p cassandra -e "ALTER USER cassandra WITH PASSWORD 'ludinnento'";