import {
  createAction, handleActions
}
from 'redux-actions'

// ------------------------------------
// Constants
// ------------------------------------
export const SET_ALL_PROJECTS = 'SET_ALL_PROJECTS'
export const DIALOG_OPEN = 'DIALOG_OPEN'
export const DIALOG_CLOSE = 'DIALOG_CLOSE'
export const CREATE_PROJECT_SUCCESS = 'CREATE_PROJECT_SUCCESS'
export const HANDLE_INPUT_CHANGE = 'HANDLE_INPUT_CHANGE'

// ------------------------------------
// Actions
// ------------------------------------
export const dialogOpen = createAction(DIALOG_OPEN, () => true)
export const dialogClose = createAction(DIALOG_CLOSE, () => false)
export const setProjects = createAction(SET_ALL_PROJECTS, (projects = []) => projects)
export const createProjectSuccess = createAction(CREATE_PROJECT_SUCCESS, () => true)
export const handleInputChange = createAction(HANDLE_INPUT_CHANGE, (e) => e.target.value)

export const getAllProjects = () => {
  return (dispatch, getState) => {
    return fetch('/api/projects')
      .then(response => {
        return response.json()
      })
      .then(json => {
        return dispatch(setProjects(json))
      })
  }
}

export const createProject = () => {
  return (dispatch, getState) => {
    return fetch('/api/project',
      {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: getState().projectName })
      })
      .then(response => {
        return response.json()
      })
      .then(json => {
        return dispatch(createProjectSuccess())
      })
  }
}

export const actions = {
  setProjects,
  getAllProjects,
  dialogOpen,
  dialogClose,
  createProject,
  handleInputChange
}

// ------------------------------------
// Reducer
// ------------------------------------
export default handleActions(
  {
    [SET_ALL_PROJECTS]: (state, {
      payload
    }) => {
      return Object.assign({}, state, {
        projects: payload
      })
    },
    [DIALOG_OPEN]: (state, { payload }) => {
      return Object.assign({}, state, {
        open: payload
      })
    },
    [DIALOG_CLOSE]: (state, { payload }) => {
      return Object.assign({}, state, {
        open: payload
      })
    },
    [CREATE_PROJECT_SUCCESS]: (state, { payload }) => {
      return Object.assign({}, state, {
        errorText: payload ? '' : 'This field is required'
      })
    }
  },

  // init state
  {
    projects: [],
    projectName: '1111',
    errorText: '',
    open: false
  }
)
