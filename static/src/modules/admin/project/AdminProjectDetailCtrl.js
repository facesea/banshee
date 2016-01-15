/*@ngInject*/
module.exports = function ($scope, $modal, $stateParams, Project) {

  $scope.loadData = function () {
    Project.getRulesByProjectId({id: $stateParams.id}).$promise
      .then(function (res) {
        $scope.rules = res;
      });
  }

  $scope.openModal = function (opt) {
    var url, ctrl;
    if (opt === 'add') {
      url = 'modules/admin/project/ruleModal.html';
      ctrl = 'RuleModalCtrl';
    }

    $modal.open({
      templateUrl: url,
      controller: ctrl
    }).result.then(function (project) {
      // TODO
    })
  }

  $scope.loadData();

}
