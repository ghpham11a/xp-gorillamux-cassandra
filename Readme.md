# Other Commands

Initialize environment

```sh
go mod init xp-gorillamux-cassandra
```

Build

```sh
go build
```

Run server

```sh
go run .
```

Run a container from the image

```sh
docker run -p 8080:8080 xp-gorillamux-cassandra
```

# 1. Build image

```
docker build -t xp-gorillamux-cassandra .
```

# 2. Setup Cassandra

```sh
helm install xp-cassandra oci://registry-1.docker.io/bitnamicharts/cassandra
```

Place Cassandra host value below to dev-accounts-configmap.yaml

```
...
Cassandra can be accessed through the following URLs from within the cluster:

  - CQL: xp-cassandra.default.svc.cluster.local:9042
...
```

Get Cassandra password. Encode it (base64) and place it in dev-accounts.secrets.yaml [YMvQ0nSpbQ]

```
export CASSANDRA_PASSWORD=$(kubectl get secret --namespace "default" xp-cassandra -o jsonpath="{.data.cassandra-password}" | base64 -d)
```

Check cluster status

```
kubectl exec -it --namespace default $(kubectl get pods --namespace default -l app.kubernetes.io/name=cassandra,app.kubernetes.io/instance=xp-cassandra -o jsonpath='{.items[0].metadata.name}') nodetool status
```

Cnnect to your Cassandra cluster using CQL:

```
kubectl run --namespace default xp-cassandra-client --rm --tty -i --restart='Never' \
   --env CASSANDRA_PASSWORD=$CASSANDRA_PASSWORD \
   --image docker.io/bitnami/cassandra:5.0.3-debian-12-r3 -- bash
```

Connect using the cqlsh client:

```
cqlsh -u cassandra -p $CASSANDRA_PASSWORD xp-cassandra
```

Create keyspace

```
cqlsh> CREATE KEYSPACE appdb WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : 1 };
```

Use it

```
USE appdb;
```

Create table

```
cassandra@cqlsh:appdb> CREATE TABLE Accounts (
    id uuid,             
    email text,          
    date_of_birth date,   
    account_number text,  
    balance decimal,      
    created_at timestamp, 
    PRIMARY KEY (id)   
);
```

Insert

```
cassandra@cqlsh:appdb> INSERT INTO accounts (id, email, date_of_birth, account_number, balance, created_at)
VALUES (68e410ae-fc9e-11ec-b939-0242ac120002, 'john.doe@example.com', '1990-01-15',
        'ACC-12345', 100.50, '2023-07-26 12:34:56')
IF NOT EXISTS;
```

Read

```
cqlsh> SELECT id, email, date_of_birth, account_number, balance, created_at FROM accounts WHERE ID = ? LIMIT 1"
```

###### 3. Setup Kafka

Install Kafka Helm chart

```sh
helm install xp-kafka oci://registry-1.docker.io/bitnamicharts/kafka
```

To check that pods were spun up

```sh
kubectl get pods --selector app.kubernetes.io/instance=xp-kafka
```

To create the topics, we will start another pod that goes into the Kafka pods and creates the topics and shuts down. To do this, we need to find the Kafka password and update it's value in the dev-kafka-topics-admin.yaml. 

```sh
# This gets an encoded value
kubectl get secret xp-kafka-user-passwords -o jsonpath='{.data.client-passwords}'

# we need the decoded version to put into the yaml
# Powershell
[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String((kubectl get secret xp-kafka-user-passwords -o jsonpath='{.data.client-passwords}')))

# bash
kubectl get secret xp-kafka-user-passwords -o jsonpath='{.data.client-passwords}' | base64 -d
```

Get the password and update the password field in dev-kafka-topics-admin.yaml

```
... required username=\"user1\" password=\"bFHKRJA2Y5\";" >> /tmp/kafka-client.properties ...
```

Then just run this command which applies the yaml updates

```sh
kubectl apply -f dev-kafka-topics-admin.yaml
```

# 3. Setup Redis

Install Redis Helm chart

```sh
helm install xp-redis oci://registry-1.docker.io/bitnamicharts/redis
```

Optional: connect to Redis

```sh
# Store the password in Powershell
$REDIS_PASSWORD = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String((kubectl get secret --namespace default xp-flask-postgres-redis -o jsonpath="{.data.redis-password}")))

# Store the password in Bash
export REDIS_PASSWORD=$(kubectl get secret --namespace default xp-flask-postgres-redis -o jsonpath="{.data.redis-password}" | base64 -d)

kubectl run --namespace default redis-client --restart='Never' --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:7.4.2-debian-12-r4 --command -- sleep infinity

kubectl exec --tty -i redis-client --namespace default -- bash

redis-cli -h xp-redis-master -p 6379

AUTH [REDIS_PASSWORD]
```