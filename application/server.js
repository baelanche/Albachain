// ExpressJS Setup
const express = require('express');
const app = express();

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 4000;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

// main page routing
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index.html');
})

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
app.post('/worker', async(req, res)=>{
    const email = req.body.email;
    console.log("add worker email: " + email);

    result = cc_call('addWorker', email)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

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

});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
