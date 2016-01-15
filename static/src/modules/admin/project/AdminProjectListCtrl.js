/*@ngInject*/
module.exports = function ($scope, $modal, Project) {
  /**
   * 加载数据
   * @param
   */
  $scope.loadData = function () {
    Project.getAllProjects().$promise
      .then(function (res) {
        $scope.projects = res;
      });
  }

  /**
   * 打开创建弹框
   * @param
   */
  $scope.openModal = function () {
    $modal.open({
      templateUrl: 'modules/admin/project/projectModal.html',
      controller: 'ProjectModalCtrl'
    }).result.then(function (project) {
      $scope.projects.push(project);
    })
  }

  $scope.loadData();

}
