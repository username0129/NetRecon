docker build -t netrecon .
docker run -d -p 80:80 -p 8081:8081 netrecon