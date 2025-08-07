# Usa una imagen oficial de Go como base
FROM golang:1.22-alpine

# Define el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum (go.sum no existe aún)
COPY go.mod go.sum ./

# Genera go.sum y descarga las dependencias. Esto es lo que soluciona el problema.
RUN go mod tidy

# Copia el resto de tu código
COPY . .

# Construye la aplicación
RUN go build -o app ./main.go

# Expone el puerto 8080 para la aplicación
EXPOSE 8080

# Define el comando para ejecutar la aplicación
CMD ["./app"]
