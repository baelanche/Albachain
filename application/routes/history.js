const express = require('express');
const router = express.Router();

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname,'../', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

router.get('/:id', function(req, res) {
    res.render('history', {wpid : req.params.id});
})

router.get('/:id/add', function(req, res) {
    res.render('historyAdd', {wpid : req.params.id});
})

router.post('/:id', async(req, res) => {
    const whid = req.body.HistoryNumber;
    const id = res.locals.user.id;
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

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
    const result = await contract.evaluateTransaction('getAllWorkHistory', whid);
    if (result == "" || result == undefined || result == null || result == 'undefined' || result == 'null') {
        res.status(200).json(result)
    } else {
        const myobj = JSON.parse(result)
        res.status(200).json(myobj)
    }
})

router.post('/:id/add', async(req, res) => {
    const id = res.locals.user.id;
    const hsid = req.body.wpid;
    const wpid = hsid.substring(0, 9);
    const yyyy = String(req.body.yyyy);
    const mm = String(req.body.mm);
    const dd = String(req.body.dd);
    const st = String(req.body.st);
    const et = String(req.body.et);
    const sttime = '' + yyyy + mm + dd + st;
    const ettime = '' + yyyy + mm + dd + et;
    var now = new Date();
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);

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
        await contract.submitTransaction("addWorkHistory", hsid, id, res.locals.user.name, wpid, "", sttime, ettime, "8000", now.toString());
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }

    res.redirect(`/history/${wpid}`);
})

module.exports = router;