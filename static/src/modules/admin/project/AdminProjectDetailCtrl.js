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

  $scope.openModal = function (event, opt) {
    $mdDialog.show({
        controller: 'RuleModalCtrl',
        templateUrl: 'modules/admin/project/ruleModal.html',
        parent: angular.element(document.body),
        targetEvent: event,
        clickOutsideToClose: true,
        fullscreen: true
      })
      .then(function (rule) {
        $scope.rules.push(rule);
      });
  };

  $scope.loadData();

};
