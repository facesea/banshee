/*@ngInject*/
module.exports = function ($scope, $stateParams, toastr, User) {
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

  $scope.loadData();

};
