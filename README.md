📌 Requisitos
1. Golang
   sudo apt update
   sudo apt install golang -y
2. Docker
   sudo apt update
   sudo apt install docker.io -y
   sudo systemctl enable --now docker
   sudo usermod -aG docker $USER
  ⚠️ Cierra sesión y vuelve a iniciar para aplicar los cambios de grupo.
3. SQLC
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ls $HOME/go/bin/sqlc
   echo 'export PATH=$PATH:/home/lucas/go/bin' >> ~/.bashrc
   source ~/.bashrc

🚀 Instalación y ejecución del proyecto
1️⃣ Clonar el proyecto
  git clone https://github.com/Smigol297/TP2-Videla-Rojas
  cd TP2-Videla-Rojas
2️⃣ Iniciar contenedores Docker
  docker compose up -d       # ejecuta en segundo plano
  docker compose up           # ejecuta en primer plano y muestra logs
3️⃣ Instalar dependencias de Go
  go mod tidy
4️⃣ Generar código Go desde SQL
  sqlc generate
5️⃣ Compilar la aplicación
  go build
6️⃣ Ejecutar la aplicación
  ./tp2

🛑 Detener contenedores
  docker compose down        # detiene contenedores
  docker compose down -v     # detiene contenedores y borra volúmenes
