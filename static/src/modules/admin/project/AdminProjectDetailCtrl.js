/*@ngInject*/
module.exports = function ($scope, $mdDialog, $state, $stateParams, toastr, Project, Rule, User, Config, Util) {
  var projectId = $scope.projectId = $stateParams.id;
  var allUsers = [];

  $scope.loadData = function () {

    // get project
    Project.get({
        id: $stateParams.id
      }).$promise
      .then(function (res) {
        $scope.project = res;
      });

    // get rules of project
    Project.getRulesByProjectId({
        id: projectId
      }).$promise
      .then(function (res) {
        $scope.rules = res;
      });

    // get users of project
    Project.getUsersByProjectId({
        id: projectId
      }).$promise
      .then(function (res) {
        $scope.users = res;
      });

    // get all users
    User.getAllUsers().$promise
      .then(function (res) {
        allUsers = res;
      });

    // get config
    Config.get().$promise
      .then(function (res) {
        $scope.config = res;
      });
  };

  $scope.edit = function() {
    Project.edit($scope.project).$promise
      .then(function() {
        toastr.success('Save Success');
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };

  $scope.deleteRule = function (event, ruleId, index) {
    var confirm = $mdDialog.confirm()
      .title('Delete Rule')
      .textContent('Would you like to delete this rule? Id: ' + ruleId)
      .ariaLabel('Delete Rule')
      .targetEvent(event)
      .ok('Yes')
      .cancel('No');
    $mdDialog.show(confirm).then(function () {
      Rule.delete({
          id: ruleId
        }).$promise
        .then(function () {
          $scope.rules.splice(index, 1);
          toastr.success('Rule Deleted!');
        })
        .catch(function (err) {
          toastr.error(err.msg);
        });
    });
  };

  $scope.deleteUser = function (event, userId, index) {
    var confirm = $mdDialog.confirm()
      .title('Remove User')
      .textContent('Would you like to remove this user from project?')
      .ariaLabel('Remove User')
      .targetEvent(event)
      .ok('Yes')
      .cancel('No');
    $mdDialog.show(confirm).then(function () {
      Project.deleteUserFromProject({
          id: projectId,
          userId: userId
        }).$promise
        .then(function () {
          $scope.users.splice(index, 1);
          toastr.success('User Deleted!');
        })
        .catch(function (err) {
          toastr.error(err.msg);
        });
    });
  };

  $scope.deleteProject = function (event) {
    var confirm = $mdDialog.confirm()
      .title('Delete Project')
      .textContent('This will delete this project with all its rules, do you want to proceed?')
      .ariaLabel('Remove Project')
      .targetEvent(event)
      .ok('Yes')
      .cancel('No');

    $mdDialog.show(confirm).then(function () {
      Project.delete({id: $scope.project.id}).$promise
        .then(function () {
          toastr.success('Project Deleted!');
          $state.go('banshee.admin.project');
        })
        .catch(function (err) {
          toastr.error(err.msg);
        });
    });

  };

  $scope.editRule = function (event, rule){
    $mdDialog.show({
        controller: 'RuleModalCtrl',
        templateUrl: 'modules/admin/project/ruleModal.html',
        parent: angular.element(document.body),
        targetEvent: event,
        clickOutsideToClose: true,
        fullscreen: true,
        bindToController: true,
        locals: {
          rule: rule,
        }
      });
  };

  $scope.openModal = function (event, opt, project) {
    var ctrl, template, users;

    if (opt === 'addRule') {
      ctrl = 'RuleModalCtrl';
      template = 'modules/admin/project/ruleModal.html';
    }

    if (opt === 'editProject') {
      ctrl = 'ProjectModalCtrl';
      template = 'modules/admin/project/projectModal.html';
    }

    if (opt === 'addUserToProject') {
      ctrl = 'UserModalCtrl';
      template = 'modules/admin/project/userModal.html';
      users = filterUsers();
    }

    $mdDialog.show({
        controller: ctrl,
        templateUrl: template,
        parent: angular.element(document.body),
        targetEvent: event,
        clickOutsideToClose: true,
        fullscreen: true,
        locals: {
          params: {
            opt: opt,
            obj: angular.copy(project) || '',
            users: users
          }
        }
      })
      .then(function (res) {
        if (opt === 'addRule') {
          $scope.rules.push(res);
        }
        if (opt === 'addUserToProject') {
          $scope.users.push(res);
        }
      });
  };

  $scope.loadData();

  /**
   * filter user:
   *  1.user.universal = true;
   *  2.user is not the existing user list;
   * @param
   */
  function filterUsers() {
    var usersIds = getUsersId();
    return allUsers.map(function(el) {
      if (!el.universal && usersIds.indexOf(el.id) < 0) {
        return el;
      }
    });
  }

  function getUsersId() {
    return $scope.users.map(function(el) {return el.id;});
  }

  $scope.buildRepr = Util.buildRepr;

};
