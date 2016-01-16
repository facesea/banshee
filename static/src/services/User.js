/**
 * Created by Panda on 16/1/16.
 */
/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/user/:id', {projectId: '@projectId', id: '@id'}, {
  });
};
