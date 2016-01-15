/*@ngInject*/
module.exports = function ($scope, $modal, $mdDialog, Project) {

  $scope.loadData = function () {
    Project.getAllProjects().$promise
      .then(function (res) {
        $scope.projects = res;
      });
  };

  $scope.openModal = function (event) {
    $mdDialog.show({
      controller: 'ProjectModalCtrl',
      templateUrl: 'modules/admin/project/projectModal.html',
      parent: angular.element(document.body),
      targetEvent: event,
      clickOutsideToClose:true,
      fullscreen: true
    })
    .then(function (project) {
      $scope.projects.push(project);
    });
  };

  $scope.loadData();

};
