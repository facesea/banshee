/*@ngInject*/
module.exports = function ($rootScope, $scope, $state) {
  $scope.includes = function(state) {
    return $state.includes(state);
  };
};
