/*@ngInject*/
module.exports = function ($scope, $modalInstance, toastr, Project) {
  $scope.project = {
    name: ''
  }
  /**
   * 提交创建表单
   * @param
   */
  $scope.create = function() {
    Project.save($scope.project).$promise
      .then(function(res) {
        $modalInstance.close(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      })
  }
}
