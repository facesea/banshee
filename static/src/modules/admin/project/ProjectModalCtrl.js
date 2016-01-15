/*@ngInject*/
module.exports = function ($scope, toastr, $mdDialog, Project) {
  $scope.project = {
    name: ''
  }

  $scope.cancel = function() {
    $mdDialog.cancel();
  }
  
  /**
   * 提交创建表单
   * @param
   */
  $scope.create = function() {
    Project.save($scope.project).$promise
      .then(function(res) {
        $mdDialog.hide(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      })
  }
}
