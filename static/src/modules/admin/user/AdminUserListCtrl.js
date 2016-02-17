/*@ngInject*/
module.exports = function ($scope, $modal, $mdDialog, $state, User) {
  $scope.autoComplete = {
    searchText: ''
  };

  $scope.loadData = function () {
    User.getAllUsers().$promise
      .then(function (res) {
        $scope.users = res;
      });
  };

  $scope.searchUser = function (user) {
    $state.go('banshee.admin.user.detail', {id: user.id});
  };

  $scope.openModal = function (event) {
    $mdDialog.show({
        controller: 'UserAddModalCtrl',
        templateUrl: 'modules/admin/user/userAddModal.html',
        parent: angular.element(document.body),
        targetEvent: event,
        clickOutsideToClose: true,
        fullscreen: true
      })
      .then(function (res) {
        $scope.users.push(res);
      });
  };

  $scope.loadData();

};
