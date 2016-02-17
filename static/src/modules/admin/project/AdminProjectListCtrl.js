/*@ngInject*/
module.exports = function ($scope, $modal, $mdDialog, $state, Project) {
  $scope.autoComplete = {
    searchText: ''
  };

  $scope.loadData = function () {
    Project.getAllProjects().$promise
      .then(function (res) {
        $scope.projects = res;
      });
  };

  $scope.searchProject = function (item) {
    $state.go('banshee.main', {project: item.id});
  }
  
  $scope.openModal = function (event) {
    $mdDialog.show({
        controller: 'ProjectModalCtrl',
        templateUrl: 'modules/admin/project/projectModal.html',
        parent: angular.element(document.body),
        targetEvent: event,
        clickOutsideToClose: true,
        fullscreen: true,
        locals: {
          params: {
            opt: 'addProject'
          }
        }
      })
      .then(function (project) {
        $scope.projects.push(project);
      });
  };

  $scope.loadData();

};
