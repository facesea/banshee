/*@ngInject*/
module.exports = function ($scope, Version) {
  $scope.loadData = function () {
    Version.get().$promise
    .then(function (res) {
      $scope.version = res.version;
    });
  };

  $scope.loadData();
};
