/*@ngInject*/
module.exports = function ($scope, $state) {
  $scope.navigateTo = function(to) {
    $state.go(to);
  }
};
