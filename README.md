Instructivo AI Model
Paso 1: Instalar Python


Paso 2: Instalar depenencias
pip install pandas fastapi uvicorn
pip install pandas scikit-learn xgboost openpyxl
pip install jinja2 python-multipart


Paso 3: Ejecutar en el ./ del folder
python -m uvicorn main:app --reload


----------------------------------------------------------


Instructivo WebApp

Paso 1: Instalar Golang

opción 1: Linux 
sudo snap install go --classic
	
opción 2: Windows
choco install go
	
opción 3: Sitio web https://go.dev/

confirmar instalación ejecutando
go version


Paso 2: Abrir terminal administradora e instalar dependencias de go ejecutando
go mod tidy


Paso 3: Ejecutar aplicación
go run ./cmd/main.go no_db

nota: debido a que la base de datos azure solo permite conectar las ips autorizadas y no hay forma de quitar esta funcionalidad, se debe ejecutar el programa con el argumento no_db.

Se preparó una demostración fija en el codigo.