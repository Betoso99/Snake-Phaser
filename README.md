# Snake-Phaser
 A simple game of snake using Phaser 3, Vue.js and Typescript

# Pre-requisites

## Local
- CockcroachDB v21.1.1
- go1.16.5 windows/amd64
- npm 8.11.0

## Docker
- Docker Cloud integration 20.10.14
- Docker Desktop 20.10.14

# How to Run

## Local 
You'll have to run a few commands
- ```npm start``` from the ```snake-front/src``` folder
- ```go run main.go``` from the ```snake-back/cmd``` folder
- ```cockroach start-single-node --insecure```

## Docker
You'll need just one command :)
- ```docker-compose up --build -d``` from the root folder

# Rules
Welcome to my snake game! It's a pleasure to meet you. I hope you will love to play this as much I did building it

First of all for you to understand the game you should know these things:

1.- You should eat the green apple, they will make you grow

2.- Avoid the red apples, they will make yow shorter twice as fast as the green ones

3.- Every 5 apples obstacles will appear, be careful

4.- You die once you eat more red apples than you can afford or you hit an obstacle

HAVE FUN :)
