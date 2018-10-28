FROM golang 
RUN mkdir /home/data 
ADD main.go /home/data/main.go
ADD types.go /home/data/types.go
ADD entrypoint.sh /home/data/entrypoint.sh
# RUN ulimit -n 99999
# RUN sysctl -w fs.file-max=100000

