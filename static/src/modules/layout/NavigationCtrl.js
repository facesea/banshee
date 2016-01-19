/*@ngInject*/
module.exports = function ($rootScope, $scope, $state, AdminNavList) {
  $rootScope.navList = AdminNavList;

  $scope.includes = function(state) {
    return $state.includes(state);
  };
};
