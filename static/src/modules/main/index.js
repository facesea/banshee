var app = angular.module('banshee.main', [])
  /*@ngInject*/
  .config(function ($stateProvider) {

    // State
    $stateProvider
      .state('banshee.main', {
        url: '/main?pattern&project',
        templateUrl: 'modules/main/list.html',
        controller: 'MainListCtrl'
      });
  })

// Controller
.controller('MainListCtrl', require('./MainListCtrl'));

module.exports = app.name;
