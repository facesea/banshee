/*@ngInject*/
module.exports = function ($scope, $modal, $mdDialog, User) {

  $scope.loadData = function () {
    User.getAllUsers().$promise
      .then(function (res) {
        $scope.users = res;
      });
  };

  $scope.openModal = function (event) {
    $mdDialog.show({
        controller: 'UserAddModalCtrl',
        templateUrl: 'modules/admin/user/userModal.html',
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
