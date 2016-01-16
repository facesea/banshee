angular.module('banshee', [

  'ngSanitize',
  'ngAnimate',
  'ngAria',
  'ngMessages',

  'ui.router',
  'ui.bootstrap',
  'ngMaterial',

  'toastr',

  'banshee.tpl',

  require('./constants'),
  require('./services'),
  require('./filters'),
  require('./directives'),

  require('./modules/layout'),
  require('./modules/admin/project')

]).config(require('./config'))
  .run(require('./run'));
