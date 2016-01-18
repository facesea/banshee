/*@ngInject*/
module.exports = function ($scope, $mdDialog, $state, $stateParams, toastr, Project, Rule, User) {
  var projectId = $stateParams.id;
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
  };

  $scope.deleteRule = function (event, ruleId, index) {
    var confirm = $mdDialog.confirm()
      .title('Delete Rule')
      .textContent('Would you like to delete this rule?')
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
  }

  $scope.openModal = function (event, opt, project) {
    var ctrl, template;

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
            users: allUsers
          }
        }
      })
      .then(function (res) {
        if (opt === 'addRule') {
          $scope.rules.push(res);
        }

        if (opt === 'editProject') {
          $scope.project = res;
        }

        if (opt === 'addUserToProject') {
          $scope.users.push(res);
        }
      });
  };

  $scope.loadData();

};
