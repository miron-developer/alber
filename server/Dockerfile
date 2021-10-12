###################################################################################################
#                                                                                                 #
#                                   Miron-developer                                               #
#                                      Al-Ber                                                     #
#                                                                                                 #
###################################################################################################

FROM golang:1.16

COPY . .
WORKDIR /pkg
RUN go mod download; go build -o ./cmd/zhibek cmd/main.go

LABEL description="This is Al-Ber project" \
    authors="Miron-developer" \
    contacts="https://github.com/miron-developer" \
    site="https://al-ber.kz"

CMD ["cmd/zhibek"]

EXPOSE 4430