/**
 * Created by Panda on 16/1/16.
 */
/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/rule/:id', {projectId: '@projectId', id: '@id'}, {
    save: {
      method: 'POST',
      url: '/api/project/:projectId/rule'
    },
    update: {
      method: 'POST',
      url: '/api/rule/:id'
    }
  });
};
