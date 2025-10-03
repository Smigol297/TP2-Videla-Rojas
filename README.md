# TP2 - Videla Rojas

## 📌 Requisitos

### 1. Golang

```bash
sudo apt update
sudo apt install golang -y
```

### 2. Docker

```bash
sudo apt update
sudo apt install docker.io -y
sudo systemctl enable --now docker
sudo usermod -aG docker $USER  # ⚠️ Cierra sesión y vuelve a iniciar para aplicar los cambios de grupo.
```

### 3. SQLC

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
ls $HOME/go/bin/sqlc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

## 🚀 Instalación y ejecución del proyecto

1. Clonar el proyecto

```bash
git clone https://github.com/Smigol297/TP2-Videla-Rojas
dd TP2-Videla-Rojas
```

2. Iniciar contenedores Docker

```bash
docker compose up -d       # ejecuta en segundo plano
docker compose up           # ejecuta en primer plano y muestra logs
```

3. Instalar dependencias de Go

```bash
go mod tidy
```

4. Generar código Go desde SQL

```bash
sqlc generate
```

5. Ejecutar tests

```bash
go test -v
```

## 🛑 Detener contenedores

```bash
docker compose down        # detiene contenedores
docker compose down -v     # detiene contenedores y borra volúmenes
```
