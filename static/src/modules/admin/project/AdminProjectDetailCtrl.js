/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, toastr, Project, Rule) {
  var projectId = $stateParams.id;

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
      Rule.delete({id: ruleId}).$promise
        .then(function() {
          $scope.rules.splice(index, 1);
          toastr.success('Rule Deleted!');
        })
        .catch(function (err) {
          toastr.error(err.msg);
        });
    });
  };

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
            obj: angular.copy(project)
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
      });
  };

  $scope.loadData();

};
