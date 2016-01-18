/*@ngInject*/
module.exports = function ($scope, $state, $stateParams, toastr, $mdDialog, User) {
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
        toastr.success('Save success!');
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };

  $scope.deleteUser = function() {
    var confirm = $mdDialog.confirm()
      .title('Delete Receiver')
      .textContent('This will delete this receiver with all its projects, do you want to proceed?')
      .ariaLabel('Remove User')
      .targetEvent(event)
      .ok('Yes')
      .cancel('No');
    $mdDialog.show(confirm).then(function () {
      User.delete({id: $scope.user.id}).$promise
        .then(function () {
          toastr.success('User Deleted!');
          $state.go('banshee.admin.user');
        })
        .catch(function (err) {
          toastr.error(err.msg);
        });
    });

  }

  $scope.loadData();

};
