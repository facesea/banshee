/*@ngInject*/
module.exports = function ($scope, $mdDialog, $stateParams, toastr, Project, params) {
  $scope.titles = {
    addUserToProject: 'Add User'
  };
  $scope.opt = params.opt;

  $scope.users = params.users;

  $scope.autoComplete = {
    searchText: ''
  };

  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    Project.addUserToProject({
      id: $stateParams.id,
      name: $scope.autoComplete.selectedItem.name
    }).$promise
    .then(function(res) {
      $mdDialog.hide($scope.autoComplete.selectedItem);
    })
    .catch(function(err) {
      toastr.error(err.msg);
    });
  };
};
