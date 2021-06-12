import axios from 'axios';
import Phaser, { Cameras } from 'phaser';

var config = {
    type: Phaser.AUTO,
    width: 800,
    height: 600,
    backgroundColor: '#0CB7F2',
    physics: {
        default: 'arcade',
    },
    scene: {
        preload: preload,
        create: create,
        update: update
    }
};

var game = new Phaser.Game(config);

var snake, cursor, apple, redApple, op, bar, ejeX, ejeY, long, high, 
    camera, scoreText, score = 0, cordinates, currentName

function preload() {
    this.load.setBaseURL('http://labs.phaser.io');
    this.load.image('block', 'assets/sprites/shinyball.png');
    this.load.image('bar', 'assets/sprites/healthbar.png');
    this.load.image('redApple', 'assets/sprites/apple.png');
    this.load.image('food', 'assets/sprites/apple.png');
}

function create() {
    currentName = localStorage.getItem('name');
    camera = this.cameras.main;
    cursor = this.add.image(200, 200, 'cursor').setVisible(true)
    snake = []
    apple = []
    bar = []
    redApple = []
    long = 1000
    high = 1000
    redApple[0] = this.physics.add.image(300, 200, 'redApple')
    apple[0] = this.physics.add.image(100, 150, 'food').setTint(0x00ff00)
    snake[0] = this.physics.add.image(200, 200, 'block').setTint(0x00ff00)
    scoreText = this.add.text(16, 16, 'Score: 0', { fontSize: '32px', fill: '#000' }).setScrollFactor(0);
    cordinates = this.add.text(16, 50, 'Sprite is at: (' + snake[0].x + ',' + snake[0].y + ')', { fontSize: '16px', fill: '#000' }).setScrollFactor(0);
    this.input.on('pointerdown', function (pointer) {
        this.input.mouse.requestPointerLock();
    }, this);
    this.input.on('pointermove', function (pointer) {
        if (this.input.mouse.locked){
            cursor.x += pointer.movementX;
            cursor.y += pointer.movementY;
            this.physics.moveToObject(snake[0], cursor, 200)
            for (let index = 1; index < snake.length; index++) {
                this.physics.moveToObject(snake[index], snake[index - 1], 200)
            }
        }
    }, this)
    
    this.rect = new Phaser.Geom.Rectangle(0, 0, long, high)
    Phaser.Actions.RandomRectangle(snake, this.rect)
    camera.startFollow(snake[0], false);

}

function update() {
    cordinates.setText('Sprite is at: (' + Math.round(snake[0].x) + ',' + Math.round(snake[0].y) + ')')
    this.physics.add.collider(snake[0], apple, Hit, null, this)
    this.physics.add.collider(snake[0], redApple, redHit, null, this)
    Phaser.Actions.WrapInRectangle(snake, this.rect)
    this.physics.world.collide(snake, bar)
}

function redHit(sna, food){
    if(snake.length > 2){
        removeItemFromArr(redApple, food)
        snake[snake.length - 1].destroy()
        snake[snake.length - 2].destroy()
        snake.splice((snake.length - 2), 2)
        food.disableBody(true, true)
        food.destroy()
    }
    else{
        axios.post('http://localhost:3000', {
            "Username": currentName,
            "Score": score.toString()
        });
        alert('Game Over')
        this.scene.pause()
        window.location.href = 'index.html'
    }
}

function Hit(sna, food) {
    snake.push(this.physics.add.image(sna.x - 20, sna.y - 20, 'block').setTint(0x00ff00))
    food.disableBody(true, true)
    apple[0].destroy()
    removeItemFromArr(apple, food)

    op=withOutDecimals(snake.length/10)
    if(op){
        bar.push(this.physics.add.staticImage(getRandomInt(20, long + 1), getRandomInt(20, high + 1), 'bar'));
        console.log('new')
        op = false
    }

    while(!apple[0]){
        ejeX = getRandomInt(20, long + 1)
        ejeY = getRandomInt(20, high + 1)
        randomApple(this, ejeX, ejeY, apple, 'food')
    }

    if(snake.length%3 === 0){
        ejeX = getRandomInt(20, long + 1)
        ejeY = getRandomInt(20, high + 1)
        randomApple(this, ejeX, ejeY, redApple, 'redApple')
        if(redApple.length === 4){
            redApple.splice(0,1)
            redApple[0].disableBody(true,true)
            redApple[0].destroy()
        }
    }
    score += 10;
    scoreText.setText('Score: ' + score);
}

function randomApple(scene, first, second, arr, path){
    var newApple = true
    
    for (let i = 0; i < bar.length; i++) {
        if(bar[i].x === ejeX && bar[i].y === ejeY){
            newApple = false
        }
    }

    if(newApple){
        if(path === 'food'){
            arr.push(scene.physics.add.image(first, second, path).setTint(0x00ff00));
        }
        else{
            arr.push(scene.physics.add.image(first, second, path));
        }
    }
}

function removeItemFromArr(arr, item) {
    var i = getIndex(arr, item)

    if (i !== -1) {
        arr.splice(i, 1);
    }
}

function getIndex(arr, item) {
    var i = arr.indexOf(item);
    return i
}

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min)) + min;
}

function withOutDecimals(numero){
    if (numero % 1 == 0) {
        return true
    } else {
        return false
    }
}