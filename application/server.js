const express = require('express');
const app = express();
const bodyParser = require('body-parser');
const mongoose = require('mongoose');

const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Hyperledger Bridge
// const { FileSystemWallet, Gateway } = require('fabric-network');

const passport = require('passport')
const session = require('express-session')
const flash = require('connect-flash');

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
const workerRouter = require('./routes/worker');
app.use('/', indexRouter);
app.use('/join', joinRouter);
app.use('/login', loginRouter);
app.use('/logout', logoutRouter);
app.use('/worker', workerRouter);

async function cc_call(fn_name, args){
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('albachain');

    var result;
    
    if(fn_name == 'addWorker')
        result = await contract.submitTransaction('addWorker', args);
    else if( fn_name == 'addWorkplace')
    {
        i=args[0];
        p=args[1];
        result = await contract.submitTransaction('addWorkplace', i, p);
    }
    else if(fn_name == 'getWorker')
        result = await contract.evaluateTransaction('getWorker', args);
    else
        result = 'not supported function'

    return result;
}

// create mate
/*app.post('/worker', async(req, res)=>{
    const email = req.body.email;
    console.log("add worker email: " + email);

    result = cc_call('addWorker', email)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})*/

// add score
app.post('/workplace', async(req, res)=>{
    const id = req.body.id;
    const place = req.body.place;

    var args=[id, place];
    result = cc_call('addWorkplace', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// find mate
/*
app.post('/worker/:id', async (req,res)=>{
    const id = req.body.id;
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('albachain');
    const result = await contract.evaluateTransaction('getWorker', id);
    const myobj = JSON.parse(result)
    res.status(200).json(myobj)
    // res.status(200).json(result)

});*/

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);

module.exports = app;