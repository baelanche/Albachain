const express = require('express');
const app = express();
const bodyParser = require('body-parser');
const mongoose = require('mongoose');
const passport = require('passport')
const session = require('express-session')
const flash = require('connect-flash');
const path = require('path');

const PORT = 4000;
const HOST = '0.0.0.0';

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: true
}))

// use static file
app.use(express.static(path.join(__dirname, 'public')));

app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');

app.use(session({
    secret: 'keyboard cat',
    resave: false,
    saveUninitialized: true
}));
app.use(passport.initialize())
app.use(passport.session())
app.use(flash())

// mongodb+srv://admin:<password>@mydb.lpk8j.mongodb.net/myFirstDatabase?retryWrites=true&w=majority
mongoose.connect('mongodb://localhost:27017/test', {
     useNewUrlParser: true, 
     useUnifiedTopology: true,
     useCreateIndex: true, 
    // useFindAndModify: false 
});
const db = mongoose.connection;
db.on('error', function(){console.log('MongoDB connection failed!')})
db.once('open', function(){console.log('MongoDB connection success!')})

// 로그인되어 있는지 확인하는 미들웨어 
app.use(function(req,res,next){
    // res.locals로 등록된 변수는 ejs 어디에서나 사용가능
    res.locals.isAuthenticated = req.isAuthenticated();
    res.locals.user = req.user;
    next();
});

const indexRouter = require('./routes/index');
const joinRouter = require('./routes/join');
const loginRouter = require('./routes/login');
const logoutRouter = require('./routes/logout');
const workplaceRouter = require('./routes/workplace');
const historyRouter = require('./routes/history');
app.use('/', indexRouter);
app.use('/join', joinRouter);
app.use('/login', loginRouter);
app.use('/logout', logoutRouter);
app.use('/workplace', workplaceRouter);
app.use('/history', historyRouter);

app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);

module.exports = app;