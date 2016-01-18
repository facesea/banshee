var app = angular.module('banshee.admin', [])
  /*@ngInject*/
  .config(function ($stateProvider) {

    // State
    $stateProvider
      .state('banshee.admin', {
        url: '/admin',
        template: '<ui-view></ui-view>',
        abstract: true
      })

      // Project router
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

      // User router
      .state('banshee.admin.user', {
        url: '/user',
        templateUrl: 'modules/admin/user/AdminUserList.html',
        controller: 'AdminUserListCtrl'
      })
      .state('banshee.admin.user.detail', {
        url: '/:id',
        views: {
          '@banshee': {
            templateUrl: 'modules/admin/user/AdminUserDetail.html',
            controller: 'AdminUserDetailCtrl'
          }
        }
      })

      // Config router
      .state('banshee.admin.config', {
        url: '/config',
        templateUrl: 'modules/admin/config/config.html',
        controller: 'AdminConfigCtrl'
      });
  })

// Controller
.controller('AdminProjectListCtrl', require('./project/AdminProjectListCtrl'))
  .controller('AdminProjectDetailCtrl', require('./project/AdminProjectDetailCtrl'))
  .controller('ProjectModalCtrl', require('./project/ProjectModalCtrl'))
  .controller('UserModalCtrl', require('./project/UserModalCtrl'))
  .controller('RuleModalCtrl', require('./project/RuleModalCtrl'))

  .controller('AdminUserListCtrl', require('./user/AdminUserListCtrl'))
  .controller('AdminUserDetailCtrl', require('./user/AdminUserDetailCtrl'))
  .controller('UserAddModalCtrl', require('./user/UserAddModalCtrl'))

  .controller('AdminConfigCtrl', require('./config/AdminConfigCtrl'));

module.exports = app.name;
