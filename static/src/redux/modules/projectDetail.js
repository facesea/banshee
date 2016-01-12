import {
  createAction, handleActions
}
from 'redux-actions'
import request from 'superagent'

// ------------------------------------
// Constants
// ------------------------------------
export const SET_ALL_RULES = 'SET_ALL_RULES'
export const GET_ALL_RULES = 'GET_ALL_RULES'
export const SET_PROJECT = 'SET_PROJECT'

export const INIT_STATE = {
  project: {},
  rules: []
}

// ------------------------------------
// Actions
// ------------------------------------
export const setProject = createAction(SET_PROJECT, (project) => project)
export const setAllRules = createAction(SET_ALL_RULES, (rules) => rules)

export const getProjectById = (id) => {
  return (dispatch, getState) => {
    request
      .get('/api/project/' + id)
      .end(function (err, res) {
        if (err || !res.ok) {
          console.error('get project error')
        } else {
          dispatch(setProject(res.body))
        }
      })
  }
}

export const getAllRules = (id) => {
  return (dispatch, getState) => {
    request
      .get('/api/project/' + id + '/rules')
      .end(function (err, res) {
        if (err || !res.ok) {
          console.error('get projects error')
        } else {
          dispatch(setAllRules(res.body))
        }
      })
  }
}

export const actions = {
  getProjectById,
  getAllRules
}

// ------------------------------------
// Reducer
// ------------------------------------
export default handleActions({
  [SET_PROJECT]: (state, { payload }) => {
    return Object.assign({}, state, {
      project: payload
    })
  },
  [SET_ALL_RULES]: (state, { payload }) => {
    return Object.assign({}, state, {
      rules: payload
    })
  }
},

// init state
INIT_STATE
)
