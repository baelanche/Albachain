const express = require('express');
const LocalStrategy = require('passport-local').Strategy;
const passport = require('passport');
const router = express.Router();
const Worker = require('../model/worker');

passport.use(new LocalStrategy(Worker.authenticate()));
passport.serializeUser(Worker.serializeUser());
passport.deserializeUser(Worker.deserializeUser());

router.get('/', function(req, res) {
    res.render('login');
})

// local-login이라는 strategy를 생성
passport.use('local-login', new LocalStrategy({
    usernameField: 'id',   // username으로 email을 사용하겠다고 선언
    passwordField: 'password',
    passReqToCallback: true // 인증을 수행하는 인증 함수로 HTTP request를 그대로 전달할지 여부를 결정

    // 로그인이 성공하면 done함수가 실행되고, done의 리턴값이 serializeUser함수의 인자로 들어가 세션을 저장한다 
  }, function(req, id, password, done){ 
    // 데이터베이스의 User 테이블에서 로그인하려는 email을 검색
    Worker.findOne({'id': id}, function(err, worker){
      if (err) return done(err);
      // DB상에 해당 email을 가진 유저가 없다면 에러 로그 출력
      if (!worker) {
        console.log("존재하지 않는 아이디입니다.")
        return done(null, false, req.flash('signinMessage', '존재하지 않는 아이디입니다.'));
      }
      // 비밀번호가 맞지 않다면 (validPassword 함수는 model/User.js에 정의되어 있음) 에러 로그 출력
      if (!worker.validPassword(password)) {
          console.log("비밀번호가 틀렸습니다.")
          return done(null, false, req.flash('signinMessage', '비밀번호가 틀렸습니다.'));
        }
      // done 함수는 자동으로 serializeUser를 호출해준다
      return done(null, worker); 
    });
  })); 

// post요청이 들어오면 local-login 전략을 실행한다.
router.post('/', passport.authenticate('local-login', 
  { failureRedirect: '/login', failureFlash: true }), function(req, res) {
  console.log("로그인 성공")
  res.redirect('/');
}); 

module.exports = router;