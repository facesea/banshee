/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, toastr, Rule, Util) {
  $scope.isEdit = false;

  if(this.rule){
    $scope.isEdit = true;
  }

  $scope.rule = this.rule || {};

  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    var params = angular.copy($scope.rule);
    if ($scope.isEdit) {
      Rule.update(params).$promise
        .then(function(res) {
          $mdDialog.hide(res);
        })
        .catch(function(err) {
          toastr.error(err.msg);
        });
    } else {
      params.projectId = $stateParams.id;
      Rule.save(params).$promise
        .then(function(res) {
          $mdDialog.hide(res);
        })
        .catch(function(err) {
          toastr.error(err.msg);
        });
    }
  };

  $scope.buildRepr = Util.buildRepr;
};
