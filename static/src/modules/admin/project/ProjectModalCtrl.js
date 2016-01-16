/*@ngInject*/
module.exports = function ($scope, toastr, $mdDialog, Project, params) {
  $scope.titles = {
    addProject: 'Create Project',
    editProject: 'Edit Project'
  };
  $scope.opt = params.opt;

  $scope.project = params.obj ? params.obj : {};

  $scope.cancel = function() {
    $mdDialog.cancel();
  };

  $scope.submit = function() {
    if (params.opt === 'editProject') {
      $scope.edit();
    } else {
      $scope.create();
    }
  };

  $scope.edit = function() {
    Project.edit($scope.project).$promise
      .then(function(res) {
        $mdDialog.hide(res);
        toastr.success('Save Success');
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };

  $scope.create = function() {
    Project.save($scope.project).$promise
      .then(function(res) {
        $mdDialog.hide(res);
      })
      .catch(function(err) {
        toastr.error(err.msg);
      });
  };
};
