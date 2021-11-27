const express = require('express');
const router = express.Router();

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname,'../', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

router.get('/', function(req, res) {
    res.render('workplace');
})

router.post('/', function(req, res) {
    if (res.locals.user.role == 'worker') {
        const id = req.body.id;
        const uname = req.body.uname;
        var works = req.body.place;
        const worknum = works.split('/')[0];
        const workname = works.split('/')[1];
        const wage = "8000";
        addWorker(id, uname, worknum, workname, wage);
    } else {
        var now = new Date();
        var d = now.getFullYear() + "" + (now.getMonth() + 1) + now.getDate();
        const id = req.body.id;
        const uname = req.body.uname;
        var works = req.body.place;
        const worknum = works.split('/')[0];
        const workname = works.split('/')[1];
        const wage = req.body.wage;
        addEmployer(id, uname, worknum, workname, d, wage);
    }
    res.redirect('/blank');
})

router.get('/:id', function(req, res) {
    res.render('workplaceInfo', {wpid: req.params.id});
})

router.post('/:id', async(req, res) => {
    const id = res.locals.user.id;
    const wpid = req.body.WorkplaceNumber;
    console.log(wpid);
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists(id);
    if (!userExists) {
        console.log('An identity for the user',id,'does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });
    const network = await gateway.getNetwork('mygreen');
    const contract = network.getContract('albachain');
    const result = await contract.evaluateTransaction('getWorkplace', wpid);
    if (result == "" || result == undefined || result == null || result == 'undefined' || result == 'null') {
        res.status(200).json(result)
    } else {
        const myobj = JSON.parse(result)
        res.status(200).json(myobj)
    }
})

async function addWorker(id, uname, worknum, workname, wage) {
    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(id);
        if (!userExists) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mygreen');

        // Get the contract from the network.
        const contract = network.getContract('albachain');

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')
        await contract.submitTransaction("addWorker", id, uname, worknum, workname, wage);
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

async function addEmployer(id, uname, worknum, workname, d, wage) {
    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(id);
        if (!userExists) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mygreen');

        // Get the contract from the network.
        const contract = network.getContract('albachain');

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')
        await contract.submitTransaction("addEmployer", id, uname, worknum, workname, d, wage);
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

module.exports = router;