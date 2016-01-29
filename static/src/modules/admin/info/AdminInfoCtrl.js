/*@ngInject*/
module.exports = function ($scope, Info) {
  $scope.info = {};
  $scope.infoText = '';

  $scope.loadData = function () {
    Info.get().$promise
      .then(function (res) {
        $scope.info= res;
        $scope.infoText = JSON.stringify(res, undefined, 2);
      });
  };

  $scope.loadData();
};
