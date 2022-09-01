
FROM docker.io/golang:1.18.0-buster

ADD ./test/image/install_docker.sh /install/install_docker.sh
RUN /install/install_docker.sh
ADD ./test/image/install_clang.sh /install/install_clang.sh
RUN /install/install_clang.sh

ADD ./test/image/install_libbpf.sh /install/install_libbpf.sh
RUN /install/install_libbpf.sh

ADD ./test/image/install_ko.sh /install/install_ko.sh
RUN /install/install_ko.sh

ADD ./test/image/install_kpt.sh /install/install_kpt.sh
RUN /install/install_kpt.sh

ADD ./test/image/install_credential_tool.sh /install/install_credential_tool.sh
RUN /install/install_credential_tool.sh

RUN apt-get -y install protobuf-compiler

RUN go install sigs.k8s.io/kind@v0.11.0
