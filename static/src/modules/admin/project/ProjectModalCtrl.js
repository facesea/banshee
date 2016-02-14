/*@ngInject*/
module.exports = function ($scope, toastr, $mdDialog, Project, params) {
  $scope.opt = params.opt;

  $scope.project = params.obj ? params.obj : {};

  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    $scope.create();
  };

  $scope.create = function() {
    Project.save($scope.project).$promise
      .then(function(res) {
        $mdDialog.hide(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };
};
