ARG BUILD_BASE_IMAGE
ARG BASE_IMAGE

FROM $BUILD_BASE_IMAGE AS softhsm-builder

RUN apt-get update && \
    apt-get install -y xz-utils zip unzip autoconf automake git libltdl-dev libssl-dev libtool openssl opensc wget && \
    rm -rf /var/lib/apt/lists/*

ENV SOFTHSM2_VERSION=2.6.1
ENV SOFTHSM2_SOURCES=/softhsm2

RUN git clone https://github.com/opendnssec/SoftHSMv2.git ${SOFTHSM2_SOURCES}
WORKDIR ${SOFTHSM2_SOURCES}
RUN git checkout ${SOFTHSM2_VERSION} -b ${SOFTHSM2_VERSION} \
    && sh autogen.sh \
    && ./configure --prefix=/usr/local --with-crypto-backend=openssl --enable-64bit --disable-gost \
    && make \
    && make install

FROM $BUILD_BASE_IMAGE AS signare-builder

ENV USER=adhara
ENV GROUP=adhara
ENV UID=1000
ENV GID=1000

ARG GOPROXY
ARG GOSUMDB

RUN addgroup --gid $GID --system $GROUP

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$GROUP" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

WORKDIR /signare

COPY . .

RUN make -C deployment build

FROM $BASE_IMAGE as signare

COPY --from=signare-builder /etc/passwd /etc/passwd
COPY --from=signare-builder /etc/group /etc/group
COPY --from=signare-builder /signare/deployment/bin/signare_linux_amd64 /signare/bin/signare

COPY --from=softhsm-builder /usr/local/lib/softhsm/libsofthsm2.so /usr/local/lib/softhsm/libsofthsm2.so

USER adhara:adhara

ENTRYPOINT [ "/signare/bin/signare" ]
