/*@ngInject*/
module.exports = function ($rootScope, $scope, $state, AdminNavList) {
  if ($state.includes('banshee.admin')) {
    $rootScope.navList = AdminNavList;
  }
  console.log(AdminNavList)

  $scope.includes = function(state) {
    return $state.includes(state);
  };
};
