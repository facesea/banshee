/*@ngInject*/
module.exports = function ($scope, Info) {
  $scope.info = {};

  $scope.loadData = function () {
    Info.get().$promise
      .then(function (res) {
        if (Object.keys(res).length === 0) {
          $scope.info = null;
        } else {
          // Tofixed with cost
          res.detectionCost = res.detectionCost.toFixed(3);
          $scope.info = res;
        }
      });
  };

  $scope.loadData();
};
