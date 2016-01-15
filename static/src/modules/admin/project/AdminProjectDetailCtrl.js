/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, Project) {

  $scope.loadData = function () {
    Project.getRulesByProjectId({id: $stateParams.id}).$promise
      .then(function (res) {
        $scope.rules = res;
      });
  }

  $scope.openModal = function (event) {
    $mdDialog.show({
      controller: 'RuleModalCtrl',
      templateUrl: 'modules/admin/project/ruleModal.html',
      parent: angular.element(document.body),
      targetEvent: event,
      clickOutsideToClose:true,
      fullscreen: true
    })
    .then(function (rule) {
      $scope.rules.push(rule)
    });
  }

  $scope.loadData();

}
