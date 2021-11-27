const express = require('express');
const LocalStrategy = require('passport-local').Strategy;
const passport = require('passport');
const router = express.Router();
const User = require('../model/user');

const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');

const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '../..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

passport.use(new LocalStrategy(User.authenticate()));
passport.serializeUser(User.serializeUser());
passport.deserializeUser(User.deserializeUser());

router.get('/', function(req, res) {
    res.render('join');
})

passport.use('local-join', new LocalStrategy({ // local-signup이라는 strategy 생성
    usernameField: 'id', // email을 username으로 사용하겠다고 선언
    passwordField: 'password',
    passReqToCallback: true  // request객체에 user의 데이터를 포함시킬지에 대한 여부를 결정, 유저 정보를 req.user로 접근할 수 있게 됨
    }, function (req, id, password, done) {
        User.findOne({id: id}, async function (err, user) { // 데이터베이스에서 넘겨받은 email으로 해당 유저가 있는지 검색
        if (err) return done(null);
        if (user) { // DB 상에 해당 유저가 있으면 에러 메시지 출력
          console.log("중복된 아이디입니다.")
          return done(null, false, req.flash('signupMessage', '중복된 아이디입니다.'));
        }
        // 저장할 유저 객체 생성 
        const newUser = new User();
        newUser.id = id;
        newUser.name = req.body.uname;
        newUser.role = req.body.role;
        // generateHash을 통해 비밀번호를 hash화
        // generateHash 함수는 model/User에 정의되어 있음
        newUser.password = newUser.generateHash(password); 
  
        newUser.save(function (err) { // 새 유저를 DB에 저장
          if (err) throw err;
          console.log("회원가입 성공")
          return done(null, newUser); // serializeUser에 user를 넘겨줌 
        });
        
        try {

          // Create a new file system based wallet for managing identities.
          const walletPath = path.join(process.cwd(), 'wallet');
          const wallet = new FileSystemWallet(walletPath);
          console.log(`Wallet path: ${walletPath}`);
  
          // Check to see if we've already enrolled the user.
          const userExists = await wallet.exists(id);
          if (userExists) {
              console.log('An identity for the user "user1" already exists in the wallet');
              return;
          }
  
          // Check to see if we've already enrolled the admin user.
          const adminExists = await wallet.exists('admin');
          if (!adminExists) {
              console.log('An identity for the admin user "admin" does not exist in the wallet');
              console.log('Run the enrollAdmin.js application before retrying');
              return;
          }
  
          // Create a new gateway for connecting to our peer node.
          const gateway = new Gateway();
          await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: false } });
  
          // Get the CA client object from the gateway for interacting with the CA.
          const ca = gateway.getClient().getCertificateAuthority();
          const adminIdentity = gateway.getCurrentIdentity();
  
          // Register the user, enroll the user, and import the new identity into the wallet.
          const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: id, role: 'client' }, adminIdentity);
          const enrollment = await ca.enroll({ enrollmentID: id, enrollmentSecret: secret });
          const userIdentity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
          wallet.import(id, userIdentity);
          console.log('Successfully registered and enrolled admin user',id,'and imported it into the wallet');
  
      } catch (error) {
          console.error(`Failed to register user "user1": ${error}`);
          process.exit(1);
      }


      })
  }));

  // post 요청이 들어오면 위에서 정의한 local-join 전략을 실행해준다
  router.post("/", passport.authenticate("local-join",{
      successRedirect:"/blank",  // 성공 혹은 실패 시 redirect되는 url
      failureRedirect:"/join",
      failureFlash: true // 실패 시 flash 메시지를 띄우는 설정
  })) 
      
  module.exports = router;