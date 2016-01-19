/**
 * Created by Panda on 16/1/13.
 */


var app = angular.module('banshee.constants', []);
app
  .constant('AdminNavList', require('./AdminNavList'))
  .constant('DateTimes', require('./DateTimes'));
module.exports = app.name;
