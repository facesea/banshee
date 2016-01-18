/*@ngInject*/
module.exports = function ($scope, Config) {
  $scope.config = {};
  $scope.configText = '';

  $scope.loadData = function () {
    Config.get().$promise
      .then(function (res) {
        $scope.config = res;
        $scope.configText = JSON.stringify(res, undefined, 2);
      });
  };

  $scope.loadData();
};
