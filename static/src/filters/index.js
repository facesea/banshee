/**
 * Created by Panda on 16/1/13.
 */

var app = angular.module('banshee.filters', [
]);

app.filter('isEmpty', [function() {
  return function(object) {
    return angular.equals({}, object);
  }
}]);

module.exports = app.name;
