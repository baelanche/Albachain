const express = require('express');
const LocalStrategy = require('passport-local').Strategy;
const passport = require('passport');
const router = express.Router();
const User = require('../model/user');

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '../', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

router.get('/blank', function(req, res) {
    res.render('blank');
})

router.get('/', function(req, res) {
    /*var result;
    if(res.locals.user != undefined && res.locals.user != 'undefined') {
        const id = res.locals.user.id;
        if (id != undefined && id != 'undefined') {
            result = getWorker(id);
        }
    }
    console.log(result);*/
    res.render('index');
})

router.post('/', async (req,res)=>{
    const id = req.body.WorkerId;
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists(id);
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }

    if(res.locals.user.role == 'worker') {
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('mygreen');
        const contract = network.getContract('albachain');
        const result = await contract.evaluateTransaction('getWorker', id);
        if (result == "" || result == undefined || result == null || result == 'undefined' || result == 'null') {
            res.status(200).json(result)
        } else {
            const myobj = JSON.parse(result)
            res.status(200).json(myobj)
        }
    } else {
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('mygreen');
        const contract = network.getContract('albachain');
        const result = await contract.evaluateTransaction('getEmployer', id);
        if (result == "" || result == undefined || result == null || result == 'undefined' || result == 'null') {
            res.status(200).json(result)
        } else {
            const myobj = JSON.parse(result)
            res.status(200).json(myobj)
        }
    }

});

module.exports = router;