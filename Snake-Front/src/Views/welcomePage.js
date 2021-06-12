window.axios = require('axios');
import Vue from 'vue/dist/vue';

var app = new Vue({
    el: '#app',
    data:{
        Records: 'Records',
        Title: 'Snake',
        users:[],
        username: '',
    },

    mounted: function(){
        axios.get('http://localhost:3000')
            .then(res => this.users = res.data);
    },

    methods:{
        play(){
            localStorage.setItem('name', this.username);
            window.location.href = 'snake.html'
        }
    }
})