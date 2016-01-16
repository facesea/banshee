/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, toastr, User, params) {
  $scope.titles = {
    addUserToProject: 'Add User'
  };
  $scope.opt = params.opt;

  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
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
