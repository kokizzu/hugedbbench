# bug report

## multi cluster
 > initiate cluster with command : ```--initial-cluster=pd0=http://pd0:2380,pd1=http://pd1:2380,pd2=http://pd2:2380``` will give error because the other `PD` didnt start yet.
 
    Error Example : 

    tidb-pd2-1    | 2022-05-31 11:43:33.843952 W | etcdserver: could not get cluster response from http://pd1:2380: Get "http://pd1:2380/members": dial tcp 172.27.0.4:2380: connect: connection refused




note : always clean tidb first before changing the version on docker-compose otherwise error will occur