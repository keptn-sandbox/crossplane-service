FROM alpine:3.13

ARG KUBE_VERSION=1.20.0
RUN wget -q https://storage.googleapis.com/kubernetes-release/release/v$KUBE_VERSION/bin/linux/amd64/kubectl -O /bin/kubectl && \
  chmod +x /bin/kubectl

ARG HELM_VERSION=3.6.3
RUN wget -q https://get.helm.sh/helm-v$HELM_VERSION-linux-amd64.tar.gz && \
	tar -zxvf helm-v$HELM_VERSION-linux-amd64.tar.gz && \
	mv linux-amd64/helm /bin/helm

######################################################################
# Keptn CLI
######################################################################

# ARG CLI_VERSION=0.9.2
# ARG CLI_DISTRO=linux-amd64

# RUN wget https://github.com/keptn/keptn/releases/download/${CLI_VERSION}/keptn-${CLI_VERSION}-${CLI_DISTRO}.tar.gz 
# RUN tar -xvf keptn-${CLI_VERSION}-${CLI_DISTRO}.tar.gz
# RUN chmod +x keptn-${CLI_VERSION}-${CLI_DISTRO}
# RUN mv keptn-${CLI_VERSION}-${CLI_DISTRO} /usr/local/bin/keptn
# RUN rm -rf keptn-${CLI_VERSION}-${CLI_DISTRO}.tar.gz
# RUN keptn version

RUN apk update && apk upgrade \
	&& apk add ca-certificates libc6-compat curl \
	&& update-ca-certificates \
	&& rm -rf /var/cache/apk/*

COPY setup.sh /setup.sh

RUN chmod +x /setup.sh

CMD [ "/bin/sh" ] 