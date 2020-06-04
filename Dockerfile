FROM golang:latest

RUN apt-get update && \
	apt-get install -y apt-utils openssh-server iputils-ping net-tools dnsutils tcpdump vim 
# && \#	rm -rf /var/lib/apt/lists/*

RUN mkdir /var/run/sshd && \
	echo 'root:xxxxxxx' | chpasswd && \
	sed -i 's/#PermitRootLogin .*/PermitRootLogin yes/' /etc/ssh/sshd_config && \
	sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd && \
	echo "export VISIBLE=now" >> /etc/profile

EXPOSE 22 

# INCLUDE APP
ENV WORK_DIR=/go/src/cloudsim/
ENV SERVICE_NAME=server
ENV EXAMPLE_NAME=client

ADD startup.sh /
RUN mkdir -p ${WORK_DIR}
COPY ${SERVICE_NAME}.go ${WORK_DIR}
COPY ${EXAMPLE_NAME}.go ${WORK_DIR}
RUN go get github.com/google/uuid
RUN go build -o ${WORK_DIR}${SERVICE_NAME} ${WORK_DIR}${SERVICE_NAME}.go

EXPOSE 23001 23002

CMD [ "/startup.sh" ]
