const express = require('express');
const router = express.Router();

const { FileSystemWallet, Gateway } = require('fabric-network');

const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '../..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

router.post('/', async(req, res, next) => {
    const email = req.body.email;
    console.log("add worker email: " + email);

    result = cc_call('addWorker', email)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

router.post('/:id', async(req,res) => {
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

module.exports = router;