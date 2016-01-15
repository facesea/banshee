/**
 * Created by Panda on 16/1/14.
 */
/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/project/:id', {id: '@id'}, {
    getAllProjects: {
      method: 'GET',
      url: '/api/projects',
      isArray: true
    },
    getRulesByProjectId: {
      method: 'GET',
      url: '/api/project/:id/rules',
      isArray: true
    }
  });
};
