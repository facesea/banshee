/**
 * Created by Panda on 16/1/14.
 */

var app = angular.module('banshee.admin', [])
  /*@ngInject*/
  .config(function ($stateProvider) {

    $stateProvider
      .state('banshee.admin', {
        url: '/admin',
        template: '<ui-view></ui-view>',
        abstract: true
      })
      .state('banshee.admin.project', {
        url: '/project',
        templateUrl: 'modules/admin/project/list.html',
        controller: 'AdminProjectListCtrl'
      })
  })

.controller('AdminProjectListCtrl', require('./AdminProjectListCtrl'))
.controller('ProjectModalCtrl', require('./ProjectModalCtrl'));

module.exports = app.name;
