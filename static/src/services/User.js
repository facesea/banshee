/**
 * Created by Panda on 16/1/16.
 */
/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/user/:id', {projectId: '@projectId', id: '@id'}, {
    getAllUsers: {
      method: 'GET',
      url: '/api/users',
      isArray: true
    },
    getProjectsByUserId: {
      method: 'GET',
      url: '/api/projects',
      isArray: true
    },
    edit: {
      method: 'PATCH',
      url: '/api/user/:id'
    }
  });
};
