# Master_slave_system
This project implements a master-slave database replication system that allows for distributed SQL query execution with asynchronous replication capabilities. The system consists of a master node that receives queries and multiple slave nodes that replicate the database state.

Key Features
Master Node
Listens on port 9000 for incoming connections from slaves

Executes SQL queries on the master MySQL database

Provides an interactive shell for direct query execution

Automatically creates databases if they don't exist

Handles multiple concurrent slave connections

Slave Node
Connects to master node to send queries

Supports both synchronous and asynchronous replication modes

Implements a buffering system for batch replication

Automatic retry mechanism for failed replications

Local query execution with MySQL integration

Architecture

