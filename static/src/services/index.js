/**
 * Created by Panda on 16/1/13.
 */

var app = angular.module('banshee.services', ['ngResource']);

app.config(function ($httpProvider) {
  $httpProvider.interceptors.push('httpInterceptor');
});

app
  .factory('httpInterceptor', require('./httpInterceptor'));

module.exports = app.name;
