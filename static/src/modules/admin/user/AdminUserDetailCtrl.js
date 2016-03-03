/*@ngInject*/
module.exports = function ($scope, $state, $stateParams, $translate, toastr, $mdDialog, User) {
  var userId = $stateParams.id;

  $scope.loadData = function () {
    // get user
    User.get({
        id: userId
      }).$promise
      .then(function (res) {
        $scope.user = res;
      });

    // get projects by user id
    User.getProjectsByUserId({
      id: userId
    }).$promise
    .then(function (res) {
      $scope.projects = res;
    });
  };

  $scope.edit = function() {
    User.edit($scope.user).$promise
      .then(function(res) {
        $scope.user = res;
        toastr.success($translate.instant('SAVE_SUCCESS'));
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };

  $scope.deleteUser = function(event) {
    var confirm = $mdDialog.confirm()
      .title($translate.instant('ADMIN_USER_DELETE_TITLE'))
      .textContent($translate.instant('ADMIN_USER_DELETE_WARN'))
      .ariaLabel($translate.instant('ADMIN_USER_DELETE_TEXT'))
      .targetEvent(event)
      .ok($translate.instant('YES'))
      .cancel($translate.instant('NO'));
    $mdDialog.show(confirm).then(function () {
      User.delete({id: $scope.user.id}).$promise
        .then(function () {
          toastr.success($translate.instant('DELETE_SUCCESS'));
          $state.go('banshee.admin.user');
        })
        .catch(function (err) {
          toastr.error(err.msg);
        });
    });

  };

  $scope.loadData();

};
