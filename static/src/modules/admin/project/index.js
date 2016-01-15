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
      .state('banshee.admin.project.detail', {
        url: '/:id',
        views: {
          '@banshee': {
            templateUrl: 'modules/admin/project/AdminProjectDetail.html',
            controller: 'AdminProjectDetailCtrl'
          }
        }
      })
  })

.controller('AdminProjectListCtrl', require('./AdminProjectListCtrl'))
.controller('AdminProjectDetailCtrl', require('./AdminProjectDetailCtrl'))
.controller('ProjectModalCtrl', require('./ProjectModalCtrl'))
.controller('RuleModalCtrl', require('./RuleModalCtrl'));

module.exports = app.name;
