const passportLocalMongoose = require('passport-local-mongoose');
const mongoose = require('mongoose');
const bcrypt = require('bcrypt-nodejs'); // 비밀번호 암호화를 위한 모듈

const workplaceSchema = new mongoose.Schema({
  workplaceNumber: {type: String, unique: true, required: true},
  users: [
      {
          type: mongoose.Schema.Types.ObjectId,
          ref: "user"
      }
  ],
  wage : {type: Number, required: true}
});
  
module.exports = mongoose.model('Workplace', workplaceSchema);