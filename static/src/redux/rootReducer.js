import { combineReducers } from 'redux'
import { routeReducer as router } from 'redux-simple-router'
import project from './modules/project'
import projectDetail from './modules/projectDetail'

export default combineReducers({
  project,
  projectDetail,
  router
})
