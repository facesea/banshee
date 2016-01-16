/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, toastr, Rule) {
  $scope.rule = {};

  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    // FIXME: Not Project.save
    var params = angular.copy($scope.rule);
    params.projectId = $stateParams.id;
    Rule.save(params).$promise
      .then(function(res) {
        $mdDialog.hide(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };
};
