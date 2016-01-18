/**
 * Created by Panda on 16/1/13.
 */
/*@ngInject*/
module.exports = function ($scope, $state) {
  $scope.includes = function(state) {
    return $state.includes(state);
  };
};
