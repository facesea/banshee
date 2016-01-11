import { combineReducers } from 'redux'
import { routeReducer as router } from 'redux-simple-router'
import project from './modules/project'

export default combineReducers({
  project,
  router
})
