# gocmdb-backend

## Description

frontend code go to `->` [gocmdb-front](https://github.com/gogo-fast/gocmdb-front)



## Quick start api server with docker

- Docker version

  ```shell
  $ docker -v
  Docker version 19.03.1, build 74b1e89
  ```

- Clone code

  ```shell
  git clone https://github.com/gogo-fast/gocmdb-backend.git
  ```

- Build images

  ```shell
  cd gocmdb-backend/apiserver
  docker build --no-cache -t cmdb-api-server:v0.1 -f dockerfiles/apiserver .
  docker build --no-cache -t cmdb-db:v0.1 -f dockerfiles/mysql dockerContexts/mysql/
  ```

- Run

  start mysql (`cmdb-db`), change `192.168.10.100` to `your_db_host_ip`

  ````shell
  docker run --name cmdb-db -d -p 192.168.10.100:3306:3306 cmdb-db:v0.1
  ````

  make sure `cmdb_db` start success

  ```shell
  docker ps -a | grep cmdb-db
  ```

  start `cmdb-api-server`, change `192.168.10.100` to `your_apiserver_host_ip`

  ```shell
  docker run --name cmdb-api-server -d --link cmdb-db:mariadb-server -p 192.168.10.100:8000:8000 cmdb-api-server:v0.1
  ```

- Forgot password

  Generate a new password with `genpass`

  ```shell
  $ docker exec -it cmdb-api-server /cmdb/apiserver/bin/genpass -p 123456  # 123456 is new pass
  123456:75a0fc76c396f590c5e63c9cc58d260e
  ```

  Update password

  ```shell
  docker exec -it cmdb-db /bin/bash -c '''mysql -ucmdb -p123456 -e "use gocmdb; update users set password = \"123456:75a0fc76c396f590c5e63c9cc58d260e\" where username = \"admin\" "; '''  
  ```

  

## Build agent docker image

- Docker version

  ```shell
  $ docker -v
  Docker version 19.03.1, build 74b1e89
  ```

- Clone code

  ```shell
  git clone https://github.com/gogo-fast/gocmdb-backend.git
  ```

- Build images

  ```shell
  cd gocmdb-backend/agent
  docker build --no-cache -t cmdb-agent:v0.1 -f dockerfiles/agent .
  ```

- Run

  start `cmdb-agent`, change `192.168.10.100` to `your_agent_host_ip`

  ```shell
  docker run --name cmdb-agent -d --link cmdb-api-server:go.cmdb.com -p 192.168.10.100:8800:8800 cmdb-agent:v0.1
  ```
  

### FQA

- x509: certificate signed by unknown authority

    ```shell script
    curl -o /etc/ssl/certs/cacert.pem https://curl.haxx.se/ca/cacert.pem
    ```


- 签名过期，请重新生成请求 or IllegalTimestamp
    
    do time sync on host

    ```shell script
    yum install ntp -y
    
    cat <<EOF >/etc/systemd/system/synctime.service
    [Unit]
    Description=Sync Time From Internet
    
    [Service]
    Type=oneshot
    ExecStart=/usr/sbin/ntpdate ntp1.aliyun.com
    EOF
    
    cat <<EOF > /etc/systemd/system/synctime-job.timer
    [Unit]
    Description=Timer for synctime.service
    
    [Timer]
    Unit=synctime.service   
    OnCalendar=*:*:00
    [Install]
    WantedBy=timers.target
    EOF
    
    systemctl daemon-reload
    systemctl start synctime-job.timer
    systemctl status synctime-job.timer
    systemctl enable synctime-job.timer
    ```