const express = require('express');
const LocalStrategy = require('passport-local').Strategy;
const passport = require('passport');
const router = express.Router();
const Worker = require('../model/worker');


router.get('/', function(req, res) {
    res.render('index');
})

module.exports = router;