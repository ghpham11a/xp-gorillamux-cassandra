# Commands

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

#### Run a container from the image

```sh
docker run -p 8080:8080 xp-gorillamux-cassandra
```

#### Install Cassandra via Helm

```sh
helm install xp-cassandra oci://registry-1.docker.io/bitnamicharts/cassandra

Cassandra can be accessed through the following URLs from within the cluster:

  - CQL: xp-cassandra.default.svc.cluster.local:9042

To get your password run:

   export CASSANDRA_PASSWORD=$(kubectl get secret --namespace "default" xp-cassandra -o jsonpath="{.data.cassandra-password}" | base64 -d)

Check the cluster status by running:

   kubectl exec -it --namespace default $(kubectl get pods --namespace default -l app.kubernetes.io/name=cassandra,app.kubernetes.io/instance=xp-cassandra -o jsonpath='{.items[0].metadata.name}') nodetool status

To connect to your Cassandra cluster using CQL:

1. Run a Cassandra pod that you can use as a client:

   kubectl run --namespace default xp-cassandra-client --rm --tty -i --restart='Never' \
   --env CASSANDRA_PASSWORD=$CASSANDRA_PASSWORD \
   --image docker.io/bitnami/cassandra:5.0.3-debian-12-r3 -- bash

2. Connect using the cqlsh client:

   cqlsh -u cassandra -p $CASSANDRA_PASSWORD xp-cassandra

To connect to your database from outside the cluster execute the following commands:

   kubectl port-forward --namespace default svc/xp-cassandra 9042:9042 &
   cqlsh -u cassandra -p $CASSANDRA_PASSWORD 127.0.0.1 9042
```

```sh
CASSANDRA_PASSWORD
```

```sh
# 1. Create keyspace
cqlsh> CREATE KEYSPACE employeedb WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : 1 };

# 2. Create table
cqlsh> CREATE TABLE employees (empID int, deptID int, first_name varchar, last_name varchar, PRIMARY KEY (empID, deptID));

# 3. Insert
cqlsh> INSERT INTO employees (empID, deptID, first_name, last_name) VALUES (1, 10, 'John', 'Smith');
cqlsh> INSERT INTO employees (empID, deptID, first_name, last_name) VALUES (2, 10, 'Jane', 'Doe');

# 4. Read
cqlsh> SELECT * FROM employees WHERE empID = 1;
```