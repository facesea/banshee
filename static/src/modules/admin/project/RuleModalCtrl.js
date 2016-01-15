/*@ngInject*/
module.exports = function ($scope, $mdDialog, toastr, Project) {
  $scope.project = {
    name: ''
  }

  $scope.cancel = function() {
    $mdDialog.cancel();
  }

  $scope.create = function() {
    // FIXME: Not Project.save
    Project.save($scope.project).$promise
      .then(function(res) {
        $mdDialog.hide(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      })
  }
}
