/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, toastr, User) {
  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    User.save($scope.user).$promise
      .then(function(res) {
        $mdDialog.hide(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };
};
