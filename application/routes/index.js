const express = require('express');
const LocalStrategy = require('passport-local').Strategy;
const passport = require('passport');
const router = express.Router();
const User = require('../model/user');


router.get('/', function(req, res) {
    res.render('index');
})

module.exports = router;