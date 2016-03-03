/**
 * Created by Panda on 16/1/13.
 */
/*@ngInject*/
module.exports = function ($scope, $state, $translate) {
  $scope.goMain = function() {
    $state.go('banshee.main', {project: '', pattern: ''}, {reload: true});
  };

  $scope.includes = function(state) {
    return $state.includes(state);
  };

  $scope.changeLanguage = function(key) {
    console.log($translate.use(key));
  };

  $scope.getLanguage = function() {
    return $translate.use();
  };
};
