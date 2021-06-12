window.axios = require('axios');
import Vue from 'vue';

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
            window.location.href = 'snake.html'
        }
    }
})

export default app.$data.username