# Client-Server-API
## Desafio 1 - FullCycle Pós Go Expert (Roberta Madge)


Falta:
 
- Entregar dois sistemas em Go: client.go e server.go
- Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida
- O timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms. 
- Client: contexto deverá retornar erro nos logs caso o tempo de execução seja insuficiente


## Run docker
```shell
 docker-compose up -d
 ```

## Acess MYSQL
 ```shell
 docker-compose exec mysql bash
 ```

## Login MYSQL
 ```shell
 mysql -uroot -p goexpert 
 ```
#### Password: root

## Create table:
```shell
create table exchange_rate (id varchar(255), name varchar(80), bid varchar(80), primary key (id));
```