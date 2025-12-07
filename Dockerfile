FROM golang:1.25-alpine AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cep-temperature /app/cmd/server/main.go

FROM scratch
WORKDIR /app

# Copia os certificados CA do sistema Alpine da fase de build, pq a imagem scratch não tem
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/cep-temperature .

EXPOSE 8080
CMD ["./cep-temperature"]